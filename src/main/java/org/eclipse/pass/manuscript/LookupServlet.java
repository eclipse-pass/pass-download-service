package org.eclipse.pass.manuscript;

/*
 * Starts the web service with the user
 * 
 * @author Maggie Olaya
 */

import java.io.IOException;
import javax.servlet.ServletException;
import javax.servlet.annotation.WebServlet;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

@WebServlet("/lookupservice")
public class LookupServlet extends HttpServlet {

    /**
    * Lookup service.
    */
    public void lookupService() {
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