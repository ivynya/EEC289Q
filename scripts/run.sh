./bin/tsp_solver_$1 --output=./solution_$2.txt ./dat/TSP_1000_euclidianDistance.txt
echo "" >> ./solution_$2.txt
./bin/tsp_solver_$1 --output=./solution_$2.txt ./dat/TSP_1000_randomDistance.txt
