#include "mpi.h"
#include <stdio.h>
int main(int argc, char *argv[] ) {
    int numprocs, rank, chunk_size, i,j;
    int max, mymax,rem;
    int matrix[800][800]; int vector[800];
    int local_matrix[800][800]; int result[800];
    int global_result[800];
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
        for(i=0;i<800;i++) {
            vector[i] = i;
            for(j=0;j<800;j++) {
                matrix[i][j] = i+j;
            }
        }
    }
    /* Distribute Vector */
    /* The vector is just send to everyone */
    MPI_Bcast(vector,800,MPI_INT,0,MPI_COMM_WORLD);
    /* Distribute Matrix */
    /* Assume the matrix is too big to bradcast. Send blocks of rows to each task,
    nrows/nprocs to each one */
    MPI_Scatter(matrix,800*chunk_size,MPI_INT,local_matrix,800*chunk_size,MPI_INT,
    0,MPI_COMM_WORLD);
    /*Each processor has a chunk of rows, now multiply and build a part of the solution vector
    */
    for(i=0;i<chunk_size;i++) {
        result[i]=0;
        for(j=0;j<800;j++) {
            result[i] += local_matrix[i][j]*vector[j];
        }
    }
    /*Send result back to master */
    MPI_Gather(result,chunk_size,MPI_INT,global_result,chunk_size,MPI_INT,
    0,MPI_COMM_WORLD);
    /*Display result */
    if(rank==0) {
        for(i=0;i<800;i++) {
            printf(" %d \t ",global_result[i]);
        }
    }
    MPI_Finalize();
    return 0;
}
