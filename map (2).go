// Project CSI2120/CSI2520
// Winter 2022
// Robert Laganiere, uottawa.ca

//Theeban Kumaresan
//300062377
//CSI2120 Winter 2022

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
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
	offset int
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

	for j := 0; j < N; j++ {
		for i := 0; i < N; i++ {

			DBscan(grid[i][j], MinPts, eps, i*10000000+j*1000000)
		}
	}

	// DBscan(grid[0][3], MinPts, eps, 0*10000000+3*1000000)

	// Parallel DBSCAN STEP 2.
	// Apply DBSCAN on each partition
	// ...

	// Parallel DBSCAN step 3.
	// merge clusters
	// *DO NOT PROGRAM THIS STEP

	end := time.Now()
	fmt.Printf("\nExecution time: %s of %d points\n", end.Sub(start), partitionSize)
}

// Applies DBSCAN algorithm on LabelledGPScoord points
// LabelledGPScoord: the slice of LabelledGPScoord points
// MinPts, eps: parameters for the DBSCAN algorithm
// offset: label of first cluster (also used to identify the cluster)
// returns number of clusters found
func DBscan(coords []LabelledGPScoord, MinPts int, eps float64, offset int) (nclusters int) {

	nclusters = 0
	//for each point P in database DB {
	for _, Point := range coords {
		//if label(P) ≠ undefined then continue
		if Point.Label != 0 {
			continue
		}
		//Neighbors N := RangeQuery(DB, distFunc, P, eps)
		var neighbours = RangeQuery(coords, Point, eps)
		//if |N| < minPts then {
		if len(neighbours) < MinPts {
			//label(P) := Noise
			Point.Label = -1
		} else {
			//C := C + 1
			nclusters++
			//label(P) := C
			cluster := Cluster{[]LabelledGPScoord{}, offset}
			Point.Label = nclusters
			cluster.coords = append(cluster.coords, Point)
			cluster.offset = offset + nclusters
			//SeedSet S := N \ {P}
			seedSet := make([]LabelledGPScoord, len(neighbours))
			copy(seedSet, neighbours)
			//for each point Q in S {

			for _, Q := range seedSet {

				//if label(Q) = Noise then label(Q) := C
				if Q.Label == -1 {
					Q.Label = cluster.offset
					cluster.coords = append(cluster.coords, Q)

				} else if Q.Label == 0 { //if label(Q) != undefined then continue

					//label(Q) := C
					Q.Label = cluster.offset
					cluster.coords = append(cluster.coords, Q)
					//Neighbors N := RangeQuery(DB, distFunc, Q, eps)
					QNeighbours := RangeQuery(coords, Q, eps)
					//if |N| ≥ minPts then {
					if len(QNeighbours)+1 >= MinPts {
						// := S ∪ N
						for _, iterator := range QNeighbours {
							if !(contains(seedSet, iterator)) {
								seedSet = append(seedSet, iterator)
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

func contains(seedSet []LabelledGPScoord, i LabelledGPScoord) (result bool) {
	for _, j := range seedSet {
		if j == i {
			return true
		}
	}
	return false
}

func distancecalc(p1 LabelledGPScoord, p2 LabelledGPScoord) (result float64) {
	var deltaLat = p2.GPScoord.lat - p1.GPScoord.lat
	var deltaLon = p2.GPScoord.long - p1.GPScoord.long
	result = math.Sqrt(math.Pow(deltaLat, 2) + math.Pow(deltaLon, 2))
	return result
}

func RangeQuery(coords []LabelledGPScoord, initial LabelledGPScoord, eps float64) (closePts []LabelledGPScoord) {
	closePts = make([]LabelledGPScoord, 0)
	for _, gpscord := range coords {
		if distancecalc(initial, gpscord) <= eps && gpscord != initial {
			closePts = append(closePts, gpscord)
		}
	}
	return closePts
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
