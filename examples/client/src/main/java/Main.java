import org.jboss.resteasy.client.jaxrs.ResteasyWebTarget;
import twitch.users.email.Test;

import javax.ws.rs.client.Client;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.client.WebTarget;

public class Main {

    public static void main(String... args) {

        String expected = "testemail@example.com";
        Test.UpdateEmailRequest request = Test.UpdateEmailRequest.newBuilder()
                .setUserId(4711L)
                .setNewEmail("test email @ exam?ple.com")
                .build();

        Client client = ClientBuilder.newClient();
        WebTarget wt = client.target("http://localhost:8080/twirp/");

        //Default java rest client
        Test.EmailBossClient ebc = new Test.EmailBossClient(wt);
        Test.UpdateEmailResponse rsp = ebc.updateEmail(request);
        System.out.println(expected + " == " + rsp.getCleanedEmail() + " -> " + expected.equals(rsp.getCleanedEmail()));


        //Using resteasy proxy client
        wt.register(new Test.ProtoBufMessageProvider());
        Test.EmailBoss eb = ((ResteasyWebTarget) wt).proxy(Test.EmailBoss.class);
        rsp = eb.updateEmail(request);
        System.out.println(expected + " == " + rsp.getCleanedEmail() + " -> " + expected.equals(rsp.getCleanedEmail()));
    }
}
