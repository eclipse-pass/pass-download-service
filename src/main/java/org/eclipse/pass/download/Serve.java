package org.eclipse.pass.manuscript;

/*
 * Starts the web service with the user
 * 
 * @author Maggie Olaya
 */

import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

@WebServlet("/downloadservice")
public class Serve extends HttpServlet{
    //list of command flags
    //call action


    public action(){
        //sets up http, fedora
        
        //creates unpaywall and downloadService objects

        //calls lookup_service_http.LookupServiceHandler, passes unpaywall variable
        //calls download_service_http.DownloadServiceHandler, passes downloadService variable
}