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



@WebServlet("/lookupservice")
public class LookupServlet extends HttpServlet{

    public void LookupService(){
        //sets up http, fedora
        
        Unpaywall unpaywall = new Unpaywall();
        LookupServiceHttp lookupService = new LookupServiceHttp();
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