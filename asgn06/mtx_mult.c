#include <stdio.h>
#include <stdlib.h>
#include <mpi.h>

int ** allocate_array(int ** array){
    int i = 0;

    if((array = (int **) malloc(sizeof(int *) * 800)) == NULL)
    {
        perror(NULL);
        exit(-1);
    }
    for(i = 0; i < 800; i++){
        if((array[i] = (int *)malloc(sizeof(int) * 800)) == NULL){
            perror(NULL);
            exit(-1);
        }
    }
    return array;
}

int main(int argc, char *argv[] ) {
    int numprocs, rank, chunk_size, i,j,k;
    int max, mymax,rem;
    int ** mtx1 = NULL; int ** mtx2 = NULL;
    int **local_matrix1 = NULL; int **local_matrix2 = NULL; int ** result = NULL;
    int ** global_result = NULL;
    int ** seq_result = NULL;
    double t1, t2;
    MPI_Status status;
    /* Initialize MPI */

    mtx1 = allocate_array(mtx1);
    mtx2 = allocate_array(mtx2);
    local_matrix1 = allocate_array(local_matrix1);
    local_matrix2 = allocate_array(local_matrix2);
    seq_result = allocate_array(seq_result);
    global_result = allocate_array(global_result);
    result = allocate_array(result);


    MPI_Init(&argc,&argv);
    MPI_Comm_rank( MPI_COMM_WORLD, &rank);
    MPI_Comm_size( MPI_COMM_WORLD, &numprocs);
    printf("Hello from process %d of %d \n",rank,numprocs);
    chunk_size = 800/numprocs;
    if (rank == 0) { /* Only on the root task... */
        /* Initialize Matrix and Vector */
        t1 = MPI_Wtime();
        for(i=0;i<800;i++) {
            for(j=0;j<800;j++) {
                seq_result[i][j] = 0;
                mtx1[i][j] = i+j;
                mtx2[i][j] = i*j;
            }
        }
        t2 = MPI_Wtime();
    }
    if (rank == 0) {
        for(i=0;i<800;i++) {
            for(j=0;j<800;j++) {
                for(k=0;k<800;k++) {
                    seq_result[i][j] = mtx1[i][k] * mtx2[k][j];
                }
            }
        }
        printf("Sequential result:\n");
        printf("Time: %f\n", t2 - t1);
    }
    /* Distribute Matricies */
    /* Assume the matrix is too big to bradcast. Send blocks of rows to each task,
    nrows/nprocs to each one */
    fprintf(stderr, "Computed Sequential result\n");
    t1 = MPI_Wtime();
    MPI_Scatter(mtx1,800*chunk_size,MPI_INT,local_matrix1,800*chunk_size,MPI_INT,0,MPI_COMM_WORLD);

    MPI_Scatter(mtx2,800*chunk_size,MPI_INT,local_matrix2,800*chunk_size,MPI_INT,0,MPI_COMM_WORLD);

    /*Each processor has a chunk of rows, now multiply and build a part of the solution vector
    */
    for(i=0;i<chunk_size;i++) {
        for(j=0;j<800;j++) {
            result[i][j] = 0;
            for(k=0;k<800;k++) {
                result[i][j] += mtx1[i][k] * mtx2[k][j];
            }
        }
    }
    /*Send result back to master */
    MPI_Gather(result,chunk_size,MPI_INT,global_result,chunk_size,MPI_INT, 0,MPI_COMM_WORLD);
    t2 = MPI_Wtime();
    /*Display result */
    if(rank==0) {
        printf("Concurrent result:\n");
        printf("Time: %f\n", t2 - t1);
    }
    MPI_Finalize();

    for(i = 0; i < 800; i++){
        for(j = 0; j < 800; j++){
            if(global_result[i][j] != result[i][j]){
                printf("Own result and MPI result disagree");
                return;
            }
        }
    }
        printf("Seq result and MPI result agree");
    return 0;
}
