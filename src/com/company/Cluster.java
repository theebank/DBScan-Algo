//Theeban Kumaresan 300062377
//Winter 2022 CSI2120
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
    //add a node to List<GPSCoord> Nodes
    public void addNode(GPSCoord node){
        this.Nodes.add(node);
        this.numofNodes++;
    }
    //returns average latitude of nodes in cluster as type double
    public double getavgLat(){
        double totLat = 0;


        for(GPSCoord i: this.Nodes){
            totLat += i.latitude;
        }

        return totLat/this.numofNodes;
    }
    //returns average longitude of nodes in cluster as type double
    public double getavgLon(){
        double totLon = 0;
        for(GPSCoord i: this.Nodes){
            totLon += i.longitude;
        }
        return totLon/this.numofNodes;

    }
}
