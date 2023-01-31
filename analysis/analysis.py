import pandas as pd
import plotext as plt

df = pd.read_csv('out.csv')
p90 = df[df.metric_name == 'http_req_duration'].groupby('timestamp').metric_value.agg(lambda x: x.quantile(0.90))
rps = df[df.metric_name == 'http_reqs'].groupby('timestamp').size()

plt.plot(p90.index, p90, label="p90 duration", yside="left")
plt.plot(rps.index, rps, label="Request per second", yside="right")
plt.title("p90 durations")
plt.show()
