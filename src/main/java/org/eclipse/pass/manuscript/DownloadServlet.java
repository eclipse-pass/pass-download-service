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
import javax.servlet.ServletException;
import java.io.IOException;



@WebServlet("/downloadservice")
public class DownloadServlet extends HttpServlet{
    //list of command flags
    //call action


    public void DownloadService(){
        //sets up http, fedora
        
        //creates unpaywall and downloadService objects

        //calls lookup_service_http.LookupServiceHandler, passes unpaywall variable
        //calls download_service_http.DownloadServiceHandler, passes downloadService variable
    }

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {
        //TODO: implement
            

    }

    @Override
    protected void doPost(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {
        //TODO: implement

    }
}