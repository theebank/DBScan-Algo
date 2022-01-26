package com.company;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class DBScan {
    public int minPts;
    public double eps;
    public DBScan(int min, double e){
        this.minPts = min;
        this.eps = e;
    }
    public List<GPSCoord> importCSV(String fname) throws IOException{
        /*
        ImportCSV method used to import data from the input CSV file
        Inputs:
            - String fName - path and filename of CSV input file
        Outputs:
            - List<TripRecord> trips - List of "TripRecord" objects which store: pickup Date, pickup Location,
                                                                                 Dropoff Location, trip Distance
        */
        List<GPSCoord> trips = new ArrayList<GPSCoord>();
        BufferedReader read = new BufferedReader(new FileReader(fname));
        read.readLine();
        String row;

        while((row = read.readLine())!= null){
            String[] data = row.split(",");

            trips.add(new GPSCoord(Double.valueOf(data[9]),Double.valueOf(data[8])));
        }
        read.close();
        return trips;
    }
    public double distancecalc(GPSCoord p1, GPSCoord p2){
        double deltaLat = p2.latitude-p1.latitude;
        double deltaLon = p2.longitude - p1.longitude;
        double result = Math.sqrt(Math.pow(deltaLat,2)+Math.pow(deltaLon,2));
        return result;
    }
    public List<GPSCoord> RangeQuery(List<GPSCoord> TripList, GPSCoord initial, double eps){
        List<GPSCoord> neighbours = new ArrayList<>();

        for (GPSCoord point: TripList) {


            if(distancecalc(initial,point)<=eps && point != initial){
                neighbours.add(point);

            }
        }
        return neighbours;
    }




}
