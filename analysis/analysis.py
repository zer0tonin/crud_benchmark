import pandas as pd
import plotext as plt

df = pd.read_csv('out.csv')
durations = df[df.metric_name == 'http_req_duration'].set_index('timestamp').groupby('timestamp')
p90 = durations['metric_value'].transform(lambda x: x.quantile(0.90))
print(p90)

plt.plot(p90)
plt.title("p90 durations")
plt.show()
