#include <stdio.h>
#include <stdlib.h>
#include <mpi.h>

#define MTX_SIZE 800

int ** allocate_array(int ** array, int rows, int cols){

    int *ptr;
    int i; int len;

    len = sizeof(int *) * rows + sizeof(int) * cols * rows;
    array = (int **)malloc(len);

    // ptr is now pointing to the first element in of 2D array
    ptr = (int *)(array + rows);

    // for loop to point rows pointer to appropriate location in 2D array
    for(i = 0; i < rows; i++)
        array[i] = (ptr + cols * i);

    return array;
}

int main(int argc, char *argv[] ) {
    int numprocs, rank, chunk_size, i,j,k;
    int max, mymax,rem;
    int ** mtx1 = NULL; int ** mtx2 = NULL;
    int **local_matrix1 = NULL; int **local_matrix2 = NULL; int ** result = NULL;
    int ** global_result = NULL;
    int ** seq_result = NULL;
    // int mtx1[MTX_SIZE][MTX_SIZE]; int mtx2[MTX_SIZE][MTX_SIZE];
    // int local_matrix1[MTX_SIZE][MTX_SIZE]; int local_matrix2[MTX_SIZE][MTX_SIZE]; int result[MTX_SIZE][MTX_SIZE];
    // int global_result[MTX_SIZE][MTX_SIZE];
    // int seq_result[MTX_SIZE][MTX_SIZE];
    double t1, t2;
    MPI_Status status;
    /* Initialize MPI */


    MPI_Init(&argc,&argv);
    MPI_Comm_rank( MPI_COMM_WORLD, &rank);
    MPI_Comm_size( MPI_COMM_WORLD, &numprocs);
    chunk_size = MTX_SIZE/numprocs;
    if (rank == 0) {
        printf("Chunk size: %d\n", chunk_size);
    }
    seq_result = allocate_array(seq_result, MTX_SIZE, MTX_SIZE);
    global_result = allocate_array(global_result, MTX_SIZE, MTX_SIZE);
    mtx1 = allocate_array(mtx1, MTX_SIZE, MTX_SIZE);
    mtx2 = allocate_array(mtx2, MTX_SIZE, MTX_SIZE);
    local_matrix1 = allocate_array(local_matrix1, MTX_SIZE, MTX_SIZE);
    local_matrix2 = allocate_array(local_matrix2, MTX_SIZE, MTX_SIZE);
    result = allocate_array(result, MTX_SIZE, MTX_SIZE);
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
        printf("Computed sequential calculation");
        printf("Time for sequential calculation: %f\n", t2 - t1);
    }
    /* Distribute Matricies */
    /* Assume the matrix is too big to bradcast. Send blocks of rows to each task,
    nrows/nprocs to each one */
    t1 = MPI_Wtime();

    MPI_Scatter(&(mtx1[0][0]),MTX_SIZE*chunk_size,MPI_INT,&(local_matrix1[0][0]),MTX_SIZE*chunk_size,MPI_INT,0,MPI_COMM_WORLD);
    MPI_Bcast(&(mtx2[0][0]),MTX_SIZE * MTX_SIZE,MPI_INT,0,MPI_COMM_WORLD);
    // if (rank == 0) {
    //     printf("RANK 0 BUFFER: %p\n", &local_matrix1);
    //     printf("MTX1:\n");
    //     for (i = 0; i < MTX_SIZE; i++) {
    //         for (j = 0; j < MTX_SIZE; j++) {
    //             printf("%d\t", local_matrix1[i][j]);
    //         }
    //         printf("\n");
    //     }
    // }
    // if (rank == 1) {
    //     printf("RANK 1 BUFFER: %p\n", &local_matrix1);
    //     printf("MTX1:\n");
    //     for (i = 0; i < MTX_SIZE; i++) {
    //         for (j = 0; j < MTX_SIZE; j++) {
    //             printf("%d\t", local_matrix1[i][j]);
    //         }
    //         printf("\n");
    //     }
    //     printf("MTX2:\n");
    //     for (i = 0; i < MTX_SIZE; i++) {
    //         for (j = 0; j < MTX_SIZE; j++) {
    //             printf("%d\t", mtx2[i][j]);
    //         }
    //         printf("\n");
    //     }
    // }

    /*Each processor has a chunk of rows, now multiply and build a part of the solution vector
    */
    for(i=0;i<chunk_size;i++) {
        for(j=0;j<MTX_SIZE;j++) {
            result[i][j] = 0;
            for(k=0;k<MTX_SIZE;k++) {
                result[i][j] += local_matrix1[i][k] * mtx2[k][j];
            }
        }
    }
    // if (rank == 1) {
    //     printf("Proc 1 result: \n");
    //     for (i = 0; i < MTX_SIZE; i++) {
    //         for (j = 0; j < MTX_SIZE; j++) {
    //             printf("%d\t", result[i][j]);
    //         }
    //         printf("\n");
    //     }
    // }
    /*Send result back to master */
    MPI_Gather(&(result[0][0]),MTX_SIZE * chunk_size,MPI_INT,&(global_result[0][0]),MTX_SIZE * chunk_size,MPI_INT, 0,MPI_COMM_WORLD);
    t2 = MPI_Wtime();
    /*Display result */
    // if(rank==0) {
    //     printf("Concurrent result:\n");
    //     for(i=0;i<MTX_SIZE;i++) {
    //         for(j=0;j<MTX_SIZE;j++) {
    //             printf(" %d \t ",global_result[i][j]);
    //         }
    //         printf("\n");
    //     }
    //     printf("Time: %f\n", t2 - t1);
    // }
    if (rank == 0) {
        printf("Concurrent result calculated\n");
        printf("Time for concurrent calculation: %f\n", t2 - t1)
    }
    if(rank == 0){
        for(i = 0; i < MTX_SIZE && i != -1; i++){
            for(j = 0; j < MTX_SIZE && j != -1; j++){
                if(global_result[i][j] != seq_result[i][j]){
                    printf("Seq result and MPI result disagree\n");
                    i = -2;
                    j = -2;
                }
            }
        }
        if (i != -1 && j != -1) {
            printf("Seq result and MPI result agree");
        }
    }
    MPI_Finalize();
    return 0;
}
