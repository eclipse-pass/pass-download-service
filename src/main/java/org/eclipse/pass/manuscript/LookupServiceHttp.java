package org.eclipse.pass.manuscript;

/*
 * Looks up a doi and provides the information associated with it
 * then creates a json file containing results
 * 
 * @author Maggie Olaya
 */

import javax.json.Json;
import javax.json.JsonObject;
import javax.json.JsonObjectBuilder;
import java.io.FileWriter;
import java.io.IOException;

public class LookupServiceHttp{

    public void lookupServiceHandler(Unpaywall unpaywall) {
        //gets doi
        String doi = "";

        createJson(unpaywall.lookup(doi));
    }

    //creates json file with manuscript info
    private void createJson(Manuscript[] manuscripts){
        //Creating a JSONObject object
        for(int i = 0; i < manuscripts.length; i++){
            JsonObject json = Json.createObjectBuilder().build();
        }
      JsonObject json = Json.createObjectBuilder().build();
   
        //TODO: implement
    }
}