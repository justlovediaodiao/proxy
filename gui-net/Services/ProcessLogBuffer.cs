namespace gui_net.Services;

public sealed class ProcessLogBuffer
{
    private const int MaxLines = 2_000;

    private readonly object _sync = new();
    private readonly Queue<string> _lines = new();

    public event Action? Changed;

    public void Append(string source, string message)
    {
        var line = $"[{source}] {message}";

        lock (_sync)
        {
            _lines.Enqueue(line);
            while (_lines.Count > MaxLines)
                _lines.Dequeue();
        }

        Changed?.Invoke();
    }

    public string GetSnapshot()
    {
        lock (_sync)
        {
            return string.Join(Environment.NewLine, _lines);
        }
    }
}
