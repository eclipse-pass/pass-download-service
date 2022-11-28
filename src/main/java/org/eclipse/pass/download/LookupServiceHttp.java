package org.eclipse.pass.manuscript;

/*
 * Looks up a doi and provides the information associated with it
 * then creates a json file containing results
 * 
 * @author Maggie Olaya
 */

import javax.json.Json;
import javax.json.JsonArray;
import javax.json.JsonObject;
import javax.json.JsonReader;

public class LookupServiceHttp{

    public lookupServiceHandler(unpaywall) {
        //gets doi

        createJson(unpaywall.lookup(doi))
    }

    //creates json file with manuscript info
    private File createJson(Manuscript[] manuscripts){
        //TODO: implement
    }
}