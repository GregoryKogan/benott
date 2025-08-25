import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd

# Set a professional and beautiful style for the plots.
sns.set_theme(style="whitegrid", palette="muted", font_scale=1.2)

# --- Benchmark Data ---
# This data is taken directly from the benchmark results.
random_data = {
    "N": [10, 100, 1000, 10000],
    "Time/Op (ns)": [4372, 115627, 1851251, 29025479],
}

grid_data = {
    "Segments": [20, 100, 200, 400],
    "Intersections": [100, 2500, 10000, 40000],
    "Time/Op (ns)": [55818, 1604035, 7506986, 31361110],
}

# Convert the data into pandas DataFrames for easy plotting.
df_random = pd.DataFrame(random_data)
df_grid = pd.DataFrame(grid_data)

# --- Chart 1: Performance vs. Number of Segments (n) ---
# This chart demonstrates the O(n log n) scaling.
plt.figure(figsize=(10, 6))
plot_random = sns.lineplot(
    x="N", y="Time/Op (ns)", data=df_random, marker="o", markersize=8
)

# Use a log-log scale to visualize the near-linearithmic relationship.
# O(n log n) appears as a slightly curved line that is close to straight on a log-log plot.
plot_random.set(xscale="log", yscale="log")

plt.title("Performance Scaling with Number of Segments (Random Data)")
plt.xlabel("Number of Segments (n)")
plt.ylabel("Time per Operation (ns, log scale)")
plt.grid(True, which="both", ls="--")

# Save the figure with high resolution for the README.
plt.savefig("benchmark_random.png", dpi=300, bbox_inches="tight")
print("Saved benchmark_random.png")


# --- Chart 2: Performance vs. Number of Intersections (k) ---
# This chart demonstrates the O(k log n) scaling in a high-contention scenario.
plt.figure(figsize=(10, 6))
plot_grid = sns.lineplot(
    x="Intersections", y="Time/Op (ns)", data=df_grid, marker="o", markersize=8
)

# A log-log scale is used here as well. A nearly straight line indicates a
# polynomial relationship, confirming that performance scales predictably with k.
plot_grid.set(xscale="log", yscale="log")

plt.title("Performance Scaling with Number of Intersections (Grid Data)")
plt.xlabel("Number of Intersections (k)")
plt.ylabel("Time per Operation (ns, log scale)")
plt.grid(True, which="both", ls="--")

# Save the figure.
plt.savefig("benchmark_grid.png", dpi=300, bbox_inches="tight")
print("Saved benchmark_grid.png")

# To display the plots if running interactively (optional)
# plt.show()
