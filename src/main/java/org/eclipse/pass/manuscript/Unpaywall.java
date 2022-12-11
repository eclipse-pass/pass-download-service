package org.eclipse.pass.manuscript;

/*
 * Looks up DOI info from unpaywall 
 * 
 * @author Maggie Olaya
 */

import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import java.net.URI;

public class Unpaywall{
    public Manuscript[] lookup(String doi) {
        Manuscript[] manuscripts = new Manuscript[5];

        //uses UnpaywallService to lookup doi, returns array of locations
        //parse through results, creating new manuscript object and adding the info

        return manuscripts;
    }
}