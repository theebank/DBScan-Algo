//Theeban Kumaresan 300062377
//Winter 2022 CSI2120
package com.company;

public class TripRecord {
    public String pickup_DateTime;
    public GPSCoord pickup_Location;
    public GPSCoord dropoff_Location;
    public double trip_distance;
    public TripRecord(String date,GPSCoord pickup, GPSCoord dropoff, double distance){
        this.pickup_DateTime = date;
        this.pickup_Location = pickup;
        this.dropoff_Location = dropoff;
        this.trip_distance = distance;
    }

}
