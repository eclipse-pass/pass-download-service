package org.eclipse.pass.manuscript;

/*
 * Verifies the doi and URL match
 * Downloads into fedora and returns URL of binary
 * 
 * @author Maggie Olaya
 */

public class DownloadService{

    public download(doi, uri){

        Unpaywall unpaywall = new Unpaywall();
        unpaywall.lookup(doi);

        //returns:
        //fedora object calls client.PostBinary
    }
}