using System;
using System.Collections.Generic;
using System.Text;
using System.Net;
using LibreHardwareMonitor.Hardware;

class SensorBridge
{
    static Computer computer;
    static int port = 8086;

    static void Main(string[] args)
    {
        foreach (var arg in args)
        {
            if (arg.StartsWith("--port="))
            {
                int.TryParse(arg.Substring("--port=".Length), out port);
            }
        }

        computer = new Computer
        {
            IsCpuEnabled = true,
            IsGpuEnabled = true,
            IsMemoryEnabled = true,
            IsMotherboardEnabled = true,
            IsControllerEnabled = true,
            IsStorageEnabled = true,
        };
        computer.Open();

        var listener = new HttpListener();
        listener.Prefixes.Add("http://127.0.0.1:" + port + "/");
        listener.Start();

        Console.Error.WriteLine("sensor_bridge: listening on port " + port);

        while (true)
        {
            try
            {
                var ctx = listener.GetContext();
                computer.Accept(new UpdateVisitor());
                string json = BuildJson();
                var buf = Encoding.UTF8.GetBytes(json);
                ctx.Response.ContentType = "application/json";
                ctx.Response.ContentLength64 = buf.Length;
                ctx.Response.OutputStream.Write(buf, 0, buf.Length);
                ctx.Response.OutputStream.Close();
            }
            catch (Exception e)
            {
                Console.Error.WriteLine("sensor_bridge: " + e.Message);
            }
        }
    }

    static string BuildJson()
    {
        var sb = new StringBuilder();
        sb.Append("{");

        bool firstHw = true;
        foreach (var hw in computer.Hardware)
        {
            bool hasSensors = false;
            var sensorSb = new StringBuilder();

            CollectSensors(hw.Sensors, ref hasSensors, sensorSb);
            foreach (var sub in hw.SubHardware)
            {
                CollectSensors(sub.Sensors, ref hasSensors, sensorSb);
            }

            if (hasSensors)
            {
                if (!firstHw) sb.Append(",");
                firstHw = false;

                string key = hw.HardwareType + "|" + hw.Name;
                sb.Append("\"" + JsonEscape(key) + "\"");
                sb.Append(":");
                sb.Append("[" + sensorSb.ToString() + "]");
            }
        }

        sb.Append("}");
        return sb.ToString();
    }

    static void CollectSensors(ISensor[] sensors, ref bool hasSensors, StringBuilder sb)
    {
        foreach (var s in sensors)
        {
            if (!s.Value.HasValue)
                continue;

            if (hasSensors)
                sb.Append(",");

            sb.Append("{");
            sb.Append("\"n\":\"" + JsonEscape(s.Name) + "\",");
            sb.Append("\"t\":\"" + JsonEscape(s.SensorType.ToString()) + "\",");
            sb.Append("\"v\":" + FloatStr(s.Value.Value));
            sb.Append("}");

            hasSensors = true;
        }
    }

    static string JsonEscape(string s)
    {
        if (string.IsNullOrEmpty(s)) return "";
        return s.Replace("\\", "\\\\").Replace("\"", "\\\"")
                .Replace("\n", "\\n").Replace("\r", "\\r").Replace("\t", "\\t");
    }

    static string FloatStr(float v)
    {
        return v.ToString("0.##", System.Globalization.CultureInfo.InvariantCulture);
    }
}

class UpdateVisitor : IVisitor
{
    public void VisitComputer(IComputer c) { c.Traverse(this); }
    public void VisitHardware(IHardware h)
    {
        h.Update();
        foreach (var sub in h.SubHardware)
            sub.Accept(this);
    }
    public void VisitSensor(ISensor s) { }
    public void VisitParameter(IParameter p) { }
}