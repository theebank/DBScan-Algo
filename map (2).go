// Project CSI2120/CSI2520
// Winter 2022
// Robert Laganiere, uottawa.ca

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

type GPScoord struct {
	lat  float64
	long float64
}

type LabelledGPScoord struct {
	GPScoord
	ID    int // point ID
	Label int // cluster ID
}
type Cluster struct {
	coords []LabelledGPScoord
	ID     int
}

const N int = 4
const MinPts int = 5
const eps float64 = 0.0003
const filename string = "yellow_tripdata_2009-01-15_9h_21h_clean.csv"

func main() {

	start := time.Now()

	gps, minPt, maxPt := readCSVFile(filename)
	fmt.Printf("Number of points: %d\n", len(gps))

	minPt = GPScoord{-74., 40.7}
	maxPt = GPScoord{-73.93, 40.8}

	// geographical limits
	fmt.Printf("NW:(%f , %f)\n", minPt.long, minPt.lat)
	fmt.Printf("SE:(%f , %f) \n\n", maxPt.long, maxPt.lat)

	// Parallel DBSCAN STEP 1.
	incx := (maxPt.long - minPt.long) / float64(N)
	incy := (maxPt.lat - minPt.lat) / float64(N)

	var grid [N][N][]LabelledGPScoord // a grid of GPScoord slices

	// Create the partition
	// triple loop! not very efficient, but easier to understand

	partitionSize := 0
	for j := 0; j < N; j++ {
		for i := 0; i < N; i++ {

			for _, pt := range gps {

				// is it inside the expanded grid cell
				if (pt.long >= minPt.long+float64(i)*incx-eps) && (pt.long < minPt.long+float64(i+1)*incx+eps) && (pt.lat >= minPt.lat+float64(j)*incy-eps) && (pt.lat < minPt.lat+float64(j+1)*incy+eps) {

					grid[i][j] = append(grid[i][j], pt) // add the point to this slide
					partitionSize++
				}
			}
		}
	}

	// ***
	// This is the non-concurrent procedural version
	// It should be replaced by a producer thread that produces jobs (partition to be clustered)
	// And by consumer threads that clusters partitions
	// for j := 0; j < N; j++ {
	// 	for i := 0; i < N; i++ {

	// 		DBscan(grid[i][j], MinPts, eps, i*10000000+j*1000000)
	// 	}
	// }
	// Parallel DBSCAN STEP 2.
	// Apply DBSCAN on each partition
	// ...
	jobs := make(chan [2]int, N*N)

	var mutex sync.WaitGroup
	threadcount := 16
	mutex.Add(threadcount)
	for i := 0; i <= threadcount; i++ {
		go Worker(jobs, grid, &mutex)
	}

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			var array [2]int
			array[0] = i
			array[1] = j
			jobs <- array
		}
	}
	close(jobs)
	mutex.Wait()

	// Parallel DBSCAN step 3.
	// merge clusters
	// *DO NOT PROGRAM THIS STEP

	end := time.Now()
	fmt.Printf("\nExecution time: %s of %d points\n", end.Sub(start), partitionSize)
}

//Worker / Consumer function
func Worker(jobs <-chan [2]int, coords [N][N][]LabelledGPScoord, done *sync.WaitGroup) {
	for {
		itr, more := <-jobs
		if more {
			i := itr[0]
			j := itr[1]
			DBscan(coords[i][j], MinPts, eps, i*10000000+j*1000000)
		} else {
			done.Done()
			return
		}
	}
}

// Applies DBSCAN algorithm on LabelledGPScoord points
// LabelledGPScoord: the slice of LabelledGPScoord points
// MinPts, eps: parameters for the DBSCAN algorithm
// offset: label of first cluster (also used to identify the cluster)
// returns number of clusters found
func DBscan(coords []LabelledGPScoord, MinPts int, eps float64, offset int) (nclusters int) {

	// *** fake code: to be rewritten
	time.Sleep(3)
	nclusters = 0

	//cluster := Cluster{[]LabelledGPScoord{}, 0}
	//for each Point P in database DB{
	for _, P := range coords {
		if P.Label != 0 {
			continue
		}
		//Neighbors N := RangeQuery(DB,distFunc, P, eps)
		N := RangeQuery(coords, P, eps)
		//if |N| <minPts then {
		if len(N) < MinPts {
			//label(P) := Noise
			P.Label = -1
		} else {
			//C := C+1
			nclusters = nclusters + 1
			//label(P) := C
			P.Label = nclusters + offset
			// var cluster Cluster
			// cluster.coords = append(cluster.coords, P)
			// cluster.ID = nclusters+offset
			var S []LabelledGPScoord
			copy(S, N)
			//for each point Q in S{
			for i := 0; i < len(S); i++ {
				Q := S[i]
				//if Label(Q) = Noise then Label(Q) := C
				if Q.Label == -1 {
					Q.Label = nclusters + offset
					//if label(Q) != undefined then continue
				} else if Q.Label == 0 {
					//label(Q) := C
					Q.Label = nclusters + offset
					//Neighbors N := RangeQuery(DB,distfunc,Q,eps)
					QN := RangeQuery(coords, Q, eps)
					//if |N| >= minPts then
					if len(QN) >= MinPts {
						//S := S u N
						for _, iterator := range QN {
							if !contains(S, iterator) {
								S = append(S, iterator)
							}
						}

					}
				}

			}
		}

	}

	// *** end of fake code.

	// End of DBscan function
	// Printing the result (do not remove)
	fmt.Printf("Partition %10d : [%4d,%6d]\n", offset, nclusters, len(coords))

	return nclusters
}
func RangeQuery(coords []LabelledGPScoord, Initial LabelledGPScoord, eps float64) (neighbours []LabelledGPScoord) {

	for _, point := range coords {
		if distancecalc(Initial, point) <= eps && point != Initial {
			neighbours = append(neighbours, point)
		}
	}
	return neighbours
}
func distancecalc(p1 LabelledGPScoord, p2 LabelledGPScoord) (result float64) {
	var deltaLat = p2.lat - p1.lat
	var deltaLon = p2.long - p1.long
	result = math.Sqrt(math.Pow(deltaLat, 2) + math.Pow(deltaLon, 2))
	return result
}
func contains(list []LabelledGPScoord, point LabelledGPScoord) bool {
	for _, i := range list {
		if i == point {
			return true
		}
	}
	return false
}

// reads a csv file of trip records and returns a slice of the LabelledGPScoord of the pickup locations
// and the minimum and maximum GPS coordinates
func readCSVFile(filename string) (coords []LabelledGPScoord, minPt GPScoord, maxPt GPScoord) {

	coords = make([]LabelledGPScoord, 0, 5000)

	// open csv file
	src, err := os.Open(filename)
	defer src.Close()
	if err != nil {
		panic("File not found...")
	}

	// read and skip first line
	r := csv.NewReader(src)
	record, err := r.Read()
	if err != nil {
		panic("Empty file...")
	}

	minPt.long = 1000000.
	minPt.lat = 1000000.
	maxPt.long = -1000000.
	maxPt.lat = -1000000.

	var n int = 0

	for {
		// read line
		record, err = r.Read()

		// end of file?
		if err == io.EOF {
			break
		}

		if err != nil {
			panic("Invalid file format...")
		}

		// get lattitude
		lat, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			fmt.Printf("\n%d lat=%s\n", n, record[8])
			panic("Data format error (lat)...")
		}

		// is corner point?
		if lat > maxPt.lat {
			maxPt.lat = lat
		}
		if lat < minPt.lat {
			minPt.lat = lat
		}

		// get longitude
		long, err := strconv.ParseFloat(record[9], 64)
		if err != nil {
			panic("Data format error (long)...")
		}

		// is corner point?
		if long > maxPt.long {
			maxPt.long = long
		}

		if long < minPt.long {
			minPt.long = long
		}

		// add point to the slice
		n++
		pt := GPScoord{lat, long}
		coords = append(coords, LabelledGPScoord{pt, n, 0})
	}

	return coords, minPt, maxPt
}
