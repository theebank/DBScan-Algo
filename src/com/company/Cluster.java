package com.company;

import java.util.ArrayList;
import java.util.List;

public class Cluster {
    public List<GPSCoord> Nodes;
    public int ClusterLabel;
    public double numofNodes = 0;
    public Cluster(int label){
        this.Nodes = new ArrayList<>();
        this.ClusterLabel = label;
    }
    public void addNode(GPSCoord node){
        this.Nodes.add(node);
        this.numofNodes++;
    }
    public double getavgLat(){
        double totLat = 0;


        for(GPSCoord i: this.Nodes){
            totLat += i.latitude;
        }

        return totLat/this.numofNodes;
    }
    public double getavgLon(){
        double totLon = 0;
        for(GPSCoord i: this.Nodes){
            totLon += i.longitude;
        }
        return totLon/this.numofNodes;

    }
}
