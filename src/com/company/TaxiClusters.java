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
        int clusterctr = 0;
        String fName = args[0];
        int minp = Integer.parseInt(args[1]);
        double eps = Double.parseDouble(args[2]);

        DBScan algo = new DBScan(minp, eps);
        List<Cluster> clusterList = new ArrayList<>();
        List<GPSCoord> nodes = algo.importCSV(fName);


        for(GPSCoord point: nodes){
            if(point.getClusterLabel()!=0){
                continue;
            }
            List<GPSCoord>neighbours = algo.RangeQuery(nodes,point,eps);

            if(neighbours.size()<minp){
                point.setClusterLabel(-1);

            }else{
            clusterctr= clusterctr+1;

            Cluster cluster = new Cluster(clusterctr);
            point.setClusterLabel(clusterctr);
            cluster.addNode(point);

            List<GPSCoord>Seedset = new ArrayList<>(neighbours);


            for (int i = 0; i < Seedset.size(); i++) {


                GPSCoord Q = Seedset.get(i);

                if(Q.getClusterLabel()==-1){
                    Q.setClusterLabel(clusterctr);
                    cluster.addNode(Q);
                }else if(Q.getClusterLabel()==0) {

                    Q.setClusterLabel(clusterctr);
                    cluster.addNode(Q);

                    List<GPSCoord> QNeighbours = algo.RangeQuery(nodes, Q, eps);

                    if (QNeighbours.size() >= minp) {
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
        FileOutputStream fout = new FileOutputStream("outputs.csv");
        BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(fout));
        String CSVSep = ",";
        StringBuffer sb = new StringBuffer();
        sb.append("ClusterLabel");
        sb.append(CSVSep);
        sb.append("Average Longitude");
        sb.append(CSVSep);
        sb.append("Average Latitude");
        sb.append(CSVSep);
        sb.append("Cluster Size");
        bw.write(sb.toString());
        bw.newLine();
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
