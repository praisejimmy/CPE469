#include <stdio.h>
#include <stdlib.h>
#include <mpi.h>

#define MTX_SIZE 800

int ** allocate_array(int ** array){
    int i = 0;

    if((array = (int **) calloc(MTX_SIZE, sizeof(int *))) == NULL)
    {
        perror(NULL);
        exit(-1);
    }
    for(i = 0; i < MTX_SIZE; i++){
        if((array[i] = (int *)calloc(MTX_SIZE, sizeof(int))) == NULL){
            perror(NULL);
            exit(-1);
        }
    }
    return array;
}

void transpose(int **array, int square_size) {
    int i, j, temp;
    for (i = 1; i < square_size; i++) {
        for (j = 0; j < i; j++) {
            temp = array[i][j];
            array[i][j] = array[j][i];
            array[j][i] = temp;
        }
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
    chunk_size = MTX_SIZE/numprocs;
    if (rank == 0) { /* Only on the root task... */
        /* Initialize Matrix and Vector */
        t1 = MPI_Wtime();
        for(i=0;i<MTX_SIZE;i++) {
            for(j=0;j<MTX_SIZE;j++) {
                seq_result[i][j] = 0;
                mtx1[i][j] = rand() % 16;
                mtx2[i][j] = rand() % 16;
            }
        }
        t2 = MPI_Wtime();
    }
    // if (rank == 0) {
    //     printf("Originals: \nMatrix1:\n");
    //     for(i=0;i<MTX_SIZE;i++) {
    //         for(j=0;j<MTX_SIZE;j++) {
    //             printf(" %d \t ",mtx1[i][j]);
    //         }
    //         printf("\n");
    //     }
    //     printf("Matrix2:\n");
    //     for(i=0;i<MTX_SIZE;i++) {
    //         for(j=0;j<MTX_SIZE;j++) {
    //             printf(" %d \t ",mtx2[i][j]);
    //         }
    //         printf("\n");
    //     }
    // }
    if (rank == 0) {
        for(i=0;i<MTX_SIZE;i++) {
            for(j=0;j<MTX_SIZE;j++) {
                for(k=0;k<MTX_SIZE;k++) {
                    seq_result[i][j] += mtx1[i][k] * mtx2[k][j];
                }
            }
        }
        // printf("Sequential result:\n");
        // for(i=0;i<MTX_SIZE;i++) {
        //     for(j=0;j<MTX_SIZE;j++) {
        //         printf(" %d \t ",seq_result[i][j]);
        //     }
        //     printf("\n");
        // }
        printf("Time: %f\n", t2 - t1);
    }
    /* Distribute Matricies */
    /* Assume the matrix is too big to bradcast. Send blocks of rows to each task,
    nrows/nprocs to each one */
    if (rank == 0) {
        printf("Computed sequential result\n");
    }
    t1 = MPI_Wtime();

    MPI_Scatter(mtx1,MTX_SIZE*chunk_size,MPI_INT,local_matrix1,MTX_SIZE*chunk_size,MPI_INT,0,MPI_COMM_WORLD);

    MPI_Bcast(mtx2,MTX_SIZE * MTX_SIZE,MPI_INT,0,MPI_COMM_WORLD);

    if (rank == 1) {
        printf("MTX1:\n");
        for (i = 0; i < 1; i++) {
            for (j = 0; j < MTX_SIZE; j++) {
                printf("%d\t", mtx1[i][j]);
            }
            printf("\n");
        }
        printf("MTX2:\n");
        // for (i = 0; i < MTX_SIZE; i++) {
        //     for (j = 0; j < MTX_SIZE; j++) {
        //         printf("%d\t", mtx2[i][j]);
        //     }
        //     printf("\n");
        // }
    }

    /*Each processor has a chunk of rows, now multiply and build a part of the solution vector
    */
    // for(i=0;i<chunk_size;i++) {
    //     for(j=0;j<MTX_SIZE;j++) {
    //         result[i][j] = 0;
    //         for(k=0;k<MTX_SIZE;k++) {
    //             result[i][j] += mtx1[i][k] * mtx2[k][j];
    //             if (rank == 1) {
    //                 printf("Calculated result: %d at i: %d, j: %d\n", result[i][j], i, j);
    //             }
    //         }
    //     }
    // }
    // /*Send result back to master */
    // MPI_Gather(result,chunk_size,MPI_INT,global_result,chunk_size,MPI_INT, 0,MPI_COMM_WORLD);
    // t2 = MPI_Wtime();
    // /*Display result */
    // if(rank==0) {
    //     printf("Concurrent result:\n");
    //     // for(i=0;i<MTX_SIZE;i++) {
    //     //     for(j=0;j<MTX_SIZE;j++) {
    //     //         printf(" %d \t ",global_result[i][j]);
    //     //     }
    //     //     printf("\n");
    //     // }
    //     printf("Time: %f\n", t2 - t1);
    // }
    //
    // if(rank == 0){
    //     for(i = 0; i < 10; i++){
    //         for(j = 0; j < 10; j++){
    //             if(global_result[i][j] != seq_result[i][j]){
    //                 printf("Own result and MPI result disagree i: %d j: %d\n", i, j);
    //             }
    //         }
    //     }
    //     printf("Seq result and MPI result agree");
    // }
    MPI_Finalize();
    return 0;
}
