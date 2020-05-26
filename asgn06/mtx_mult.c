#include <stdio.h>
#include <stdlib.h>
#include <mpi.h>

int ** allocate_array(int ** array){
    int i = 0;

    if((array = (int **) calloc(10, sizeof(int *))) == NULL)
    {
        perror(NULL);
        exit(-1);
    }
    for(i = 0; i < 10; i++){
        if((array[i] = (int *)calloc(10, sizeof(int))) == NULL){
            perror(NULL);
            exit(-1);
        }
    }
    return array;
}

void transpose(int **array, int square_size) {
    int i, j, temp;
    printf("Before tranpose inside function\n");
    for (i = 0; i < 10; i++) {
        for (j = 0; j < 10; j++) {
            printf("%d\t", array[i][j]);
        }
        printf("\n");
    }
    for (i = 0; i < square_size; i++) {
        for (j = i+1; j < square_size; j++) {
            temp = array[i][j];
            printf("Temp: %d array[%d][%d]: %d array[%d][%d]: %d\n", temp, i, j, array[i][j], j, i, array[j][i]);
            array[i][j] = array[j][i];
            array[j][i] = temp;
        }
    }
    printf("After tranpose in function\n");
    for (i = 0; i < 10; i++) {
        for (j = 0; j < 10; j++) {
            printf("%d\t", array[i][j]);
        }
        printf("\n");
    }
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
    chunk_size = 10/numprocs;
    if (rank == 0) { /* Only on the root task... */
        /* Initialize Matrix and Vector */
        t1 = MPI_Wtime();
        for(i=0;i<10;i++) {
            for(j=0;j<10;j++) {
                seq_result[i][j] = 0;
                mtx1[i][j] = i+j;
                mtx2[i][j] = i*j;
            }
        }
        t2 = MPI_Wtime();
    }
    if (rank == 0) {
        for(i=0;i<10;i++) {
            for(j=0;j<10;j++) {
                for(k=0;k<10;k++) {
                    seq_result[i][j] = mtx1[i][k] * mtx2[k][j];
                }
            }
        }
        printf("Sequential result:\n");
        for(i=0;i<10;i++) {
            for(j=0;j<10;j++) {
                printf(" %d \t ",seq_result[i][j]);
            }
            printf("\n");
        }
        printf("Time: %f\n", t2 - t1);
    }
    /* Distribute Matricies */
    /* Assume the matrix is too big to bradcast. Send blocks of rows to each task,
    nrows/nprocs to each one */
    fprintf(stdout, "Computed Sequential result\n");
    fflush(stdout);
    t1 = MPI_Wtime();
    if (rank == 0) {
        printf("Mtx 2 before:\n");
        for (i = 0; i < 10; i++) {
            for (j = 0; j < 10; j++) {
                printf("%d\t", mtx2[i][j]);
            }
            printf("\n");
        }
        printf("Mtx2 after transpose:\n");

        transpose(mtx2, 10);
        for (i = 0; i < 10; i++) {
            for (j = 0; j < 10; j++) {
                printf("%d\t", mtx2[i][j]);
            }
            printf("\n");
        }
    }

    MPI_Scatter(mtx1,10*chunk_size,MPI_INT,local_matrix1,10*chunk_size,MPI_INT,0,MPI_COMM_WORLD);

    MPI_Scatter(mtx2,10*chunk_size,MPI_INT,local_matrix2,10*chunk_size,MPI_INT,0,MPI_COMM_WORLD);

    /*Each processor has a chunk of rows, now multiply and build a part of the solution vector
    */
    for(i=0;i<chunk_size;i++) {
        for(j=0;j<10;j++) {
            result[i][j] = 0;
            for(k=0;k<10;k++) {
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
        for(i=0;i<10;i++) {
            for(j=0;j<10;j++) {
                printf(" %d \t ",global_result[i][j]);
            }
            printf("\n");
        }
        printf("Time: %f\n", t2 - t1);
    }

    if(rank == 0){
        for(i = 0; i < 10; i++){
            for(j = 0; j < 10; j++){
                if(global_result[i][j] != seq_result[i][j]){
                    printf("Own result and MPI result disagree i: %d j: %d\n", i, j);
                }
            }
        }
        printf("Seq result and MPI result agree");
    }
    MPI_Finalize();
    return 0;
}
