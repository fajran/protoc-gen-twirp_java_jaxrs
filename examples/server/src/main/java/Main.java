import io.undertow.Undertow;
import org.jboss.resteasy.plugins.server.undertow.UndertowJaxrsServer;
import twitch.users.email.Test;

import javax.ws.rs.ApplicationPath;
import javax.ws.rs.core.Application;
import java.util.HashSet;
import java.util.Set;

@ApplicationPath("/twirp")
public class Main extends Application {

    @Override
    public Set<Object> getSingletons() {
        HashSet set = new HashSet<>();
        set.add(new TestRpc());
        set.add(new Test.ProtoBufMessageProvider());
        return set;
    }

    public static void main(String... args) {
        UndertowJaxrsServer server = new UndertowJaxrsServer().start(Undertow.builder().addHttpListener(8080, "localhost"));
        server.deploy(new Main());
    }
}
