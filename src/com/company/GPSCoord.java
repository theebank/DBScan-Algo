//Theeban Kumaresan 300062377
//Winter 2022 CSI2120
package com.company;

public class GPSCoord {
    public double latitude;
    public double longitude;

    public int ClusterLabel;
    // 0 is unvisited
    // -1 is NOISE
    // >0 is cluster label

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
