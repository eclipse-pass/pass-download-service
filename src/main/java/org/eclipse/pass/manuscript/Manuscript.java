package org.eclipse.pass.manuscript;

/*
 * Contains the info for the manuscripts that
 * will be found during DOI lookup
 * 
 * @author Maggie Olaya
 */
public class Manuscript {
    private String location; //location URI of manuscript
    private String repoInstitution; //readable label for the repository where the article can be found
    private String type; //the MIME type of manuscript file
    private String source; //the API where we found the file
    private String name; //the file name

  /**
   * Manuscript creation.
   * 
   */
    public void manuscript(String location, String repo, String type, String source, String name) {
        this.location = location;
        repo = repoInstitution;
        this.type = type;
        this.source = source;
        this.name = name;
    }

    public String getLocation() {
        return location;
    }

    public String getRepoInstitution() {
        return repoInstitution;
    }

    public String getType() {
        return type;
    }

    public String getSource() {
        return source;
    }

    public String getName() {
        return name;
    }
}
