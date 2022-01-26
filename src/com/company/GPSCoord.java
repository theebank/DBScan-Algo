package com.company;

public class GPSCoord {
    public double latitude;
    public double longitude;
    public int ClusterLabel;


    public GPSCoord(double latitude, double longitude){
        this.latitude = latitude;
        this.longitude = longitude;
        this.ClusterLabel = 0;
    }

    public void setClusterLabel(int clusterLabel) {
        ClusterLabel = clusterLabel;
    }

    public int getClusterLabel() {
        return ClusterLabel;
    }
}
