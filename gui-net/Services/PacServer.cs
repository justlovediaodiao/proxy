using System.Net;
using System.Text;

namespace gui_net.Services;

public class PacServer
{
    private HttpListener? _listener;
    private Thread? _serverThread;
    private string? _pacContent;
    private bool _running;

    public void Start(string host, int port)
    {
        Stop();

        var pacPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "resources", "pac.js");
        if (!File.Exists(pacPath))
            throw new FileNotFoundException($"PAC file not found: {pacPath}", pacPath);

        _pacContent = File.ReadAllText(pacPath);

        _listener = new HttpListener();
        // HttpListener requires admin rights for some prefixes, but localhost usually works if not reserved.
        // Using + or * might require admin.
        _listener.Prefixes.Add($"http://{host}:{port}/");
        _listener.Start();

        _running = true;
        _serverThread = new Thread(Listen)
        {
            IsBackground = true
        };
        _serverThread.Start();
    }

    public void Stop()
    {
        _running = false;
        _listener?.Stop();
        _listener?.Close();
        _listener = null;
        _pacContent = null;
    }

    private void Listen()
    {
        while (_running && _listener != null && _listener.IsListening)
        {
            try
            {
                var context = _listener.GetContext();
                ProcessRequest(context);
            }
            catch (HttpListenerException)
            {
                // Listener stopped
            }
            catch (Exception ex)
            {
                Console.WriteLine($"PAC Server error: {ex.Message}");
            }
        }
    }

    private void ProcessRequest(HttpListenerContext context)
    {
        try
        {
            var response = context.Response;
            var pacContent = _pacContent ?? throw new InvalidOperationException("PAC content has not been loaded.");

            var buffer = Encoding.UTF8.GetBytes(pacContent);
            response.ContentLength64 = buffer.Length;
            response.ContentType = "application/x-ns-proxy-autoconfig";
            response.OutputStream.Write(buffer, 0, buffer.Length);
            response.OutputStream.Close();
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error processing request: {ex.Message}");
        }
    }
}
