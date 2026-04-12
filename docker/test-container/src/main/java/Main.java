import org.arl.fjage.Platform;
import org.arl.fjage.RealTimePlatform;
import org.arl.fjage.connectors.Connector;
import org.arl.fjage.connectors.WebServer;
import org.arl.fjage.connectors.WebSocketHubConnector;
import org.arl.fjage.remote.MasterContainer;
import org.arl.fjage.shell.ConsoleShell;
import org.arl.fjage.shell.GroovyScriptEngine;
import org.arl.fjage.shell.ShellAgent;

import java.io.InputStream;
import java.util.logging.LogManager;

public class Main {
  public static void main(String[] args) throws Exception {
    {
      final InputStream inputStream = Thread.currentThread().getContextClassLoader().getResourceAsStream("logging.properties");
      if (inputStream != null) {
        LogManager.getLogManager().readConfiguration(inputStream);
      }
    }

    final Platform platform = new RealTimePlatform();
    final MasterContainer container = new MasterContainer(platform);
    WebServer.getInstance(8080, "0.0.0.0").addStatic("/", "/org/arl/fjage/web");
    final Connector conn = new WebSocketHubConnector(8080, "/shellAgent/ws");
    container.openWebSocketServer(8080, "/ws");
    final ShellAgent shellAgent = new ShellAgent(new ConsoleShell(conn), new GroovyScriptEngine());
    container.add("shell", shellAgent);
    platform.start();
  }
}
