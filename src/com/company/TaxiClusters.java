//Theeban Kumaresan 300062377
//Winter 2022 CSI2120
package com.company;

import java.io.BufferedWriter;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.OutputStreamWriter;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class TaxiClusters {
    public static void main(String[] args) throws IOException {
        //C := 0
        int clusterctr = 0;
        //Getting arguments from running program
        String fName = args[0];
        int minp = Integer.parseInt(args[1]);
        double eps = Double.parseDouble(args[2]);
        //Initializing lists and DBScan object
        DBScan algo = new DBScan(minp, eps);
        List<Cluster> clusterList = new ArrayList<>();
        List<GPSCoord> nodes = algo.importCSV(fName);

        //DBScan algorithm
        //for each point P in database DB {
        for(GPSCoord point: nodes){
            //if label(P) ≠ undefined then continue
            if(point.getClusterLabel()!=0){
                continue;
            }
            //Neighbors N := RangeQuery(DB, distFunc, P, eps) 
            List<GPSCoord>neighbours = algo.RangeQuery(nodes,point,eps);
            //if |N| < minPts then {
            if(neighbours.size()<minp){
                //label(P) := Noise
                point.setClusterLabel(-1);

            }else{
            //C := C + 1
            clusterctr= clusterctr+1;
            //label(P) := C
            Cluster cluster = new Cluster(clusterctr);
            point.setClusterLabel(clusterctr);
            cluster.addNode(point);
            //SeedSet S := N \ {P}
            List<GPSCoord>Seedset = new ArrayList<>(neighbours);

            //for each point Q in S {
            for (int i = 0; i < Seedset.size(); i++) {
                GPSCoord Q = Seedset.get(i);
                //if label(Q) = Noise then label(Q) := C
                if(Q.getClusterLabel()==-1){
                    Q.setClusterLabel(clusterctr);
                    cluster.addNode(Q);
                //if label(Q) != undefined then continue
                }else if(Q.getClusterLabel()==0) {
                    //label(Q) := C
                    Q.setClusterLabel(clusterctr);
                    cluster.addNode(Q);
                    //Neighbors N := RangeQuery(DB, distFunc, Q, eps)
                    List<GPSCoord> QNeighbours = algo.RangeQuery(nodes, Q, eps);
                    //if |N| ≥ minPts then { 
                    if (QNeighbours.size() >= minp) {
                        // := S ∪ N
                        for(GPSCoord iterator: QNeighbours){
                            if(!(Seedset.contains(iterator))){
                                Seedset.add(iterator);
                            }
                        }

                    }
                }



            }
            clusterList.add(cluster);



            }

        }
        //Writing to CSV
        FileOutputStream fout = new FileOutputStream("outputs.csv");//Change to desired path, otherwise will be placed in "src" folder
        BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(fout));
        String CSVSep = ",";
        StringBuffer sb = new StringBuffer();
        //Headers for each column
        sb.append("ClusterLabel");
        sb.append(CSVSep);
        sb.append("Average Longitude");
        sb.append(CSVSep);
        sb.append("Average Latitude");
        sb.append(CSVSep);
        sb.append("Cluster Size");
        bw.write(sb.toString());
        bw.newLine();
        //Looping through each cluster
        for(Cluster cluster:clusterList ){
            sb = new StringBuffer();
            sb.append(cluster.ClusterLabel);
            sb.append(CSVSep);
            sb.append(cluster.getavgLon());
            sb.append(CSVSep);
            sb.append(cluster.getavgLat());
            sb.append(CSVSep);
            sb.append(cluster.Nodes.size());
            bw.write(sb.toString());
            bw.newLine();
        }
        bw.flush();
        bw.close();

    }
}
