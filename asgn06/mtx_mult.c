#include "mpi.h"
#include <stdio.h>

int main(int argc, char *argv[] ) {
    int numprocs, rank, chunk_size, i,j,k;
    int max, mymax,rem;
    int mtx1[800][800]; int mtx2[800][800];
    int local_matrix1[800][800]; int local_matrix2; int result[800][800];
    int global_result[800][800];
    int seq_result[800][800];
    double t1, t2;
    MPI_Status status;
    /* Initialize MPI */
    MPI_Init( &argc,&argv);
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
        for(i=0;i<800;i++) {
            for(j=0;j<800;j++) {
                printf(" %d \t ",global_result[i][j]);
            }
            printf("\n");
        }
        printf("Time: %f\n", t2 - t1);
    }
    /* Distribute Matricies */
    /* Assume the matrix is too big to bradcast. Send blocks of rows to each task,
    nrows/nprocs to each one */
    t1 = MPI_Wtime();
    MPI_Scatter(mtx1,800*chunk_size,MPI_INT,local_matrix1,800*chunk_size,MPI_INT,
        (void *) 0,MPI_COMM_WORLD);

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
    MPI_Gather(result,chunk_size,MPI_INT,global_result,chunk_size,MPI_INT,
    0,MPI_COMM_WORLD);
    t2 = MPI_Wtime();
    /*Display result */
    if(rank==0) {
        printf("Concurrent result:\n");
        for(i=0;i<800;i++) {
            for(j=0;j<800;j++) {
                printf(" %d \t ",global_result[i][j]);
            }
            printf("\n")
        }
        printf("Time: %f\n", t2 - t1);
    }
    MPI_Finalize();
    return 0;
}
