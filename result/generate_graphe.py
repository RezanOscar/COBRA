import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import glob
import seaborn as sns

# Data extracted from "Time for task.txt"
task_time_data = {
    'Tasks': [10, 20, 30, 40, 50, 60, 70, 80],
    'RoundRobin': [26.84, 52.36, 79.03, 105.80, 137.26, 160.05, 178.01, 206.87],
    'Random': [26.14, 52.62, 80.85, 104.24, 135.39, 159.48, 177.86, 205.09],
    'ECP': [27.35, 56.58, 84.12, 110.40, 141.04, 170.49, 193.82, 221.74],
    'EnergyAware': [28.29, 59.02, 88.95, 116.62, 144.53, 172.11, 194.6, 223.463],
    'Cobra': [28.78, 57.46, 87.03, 109.46, 141.92, 165.07, 188.27, 215.22]
}

# Data for NumPeer, Transaction time, and Consensus time
latency_data = {
    'NumPeer': [5, 10, 20, 30, 40],
    'TransactionTime': [2.87, 5.94, 12.32, 21.24, 30.35],
    'ConsensusTime': [0.7, 1.89, 4.12, 6.90, 9.24]
}

# Confidence intervals as percentages
confidence_intervals = {
    'Random': 96,
    'RoundRobin': 96.1,
    'ECP': 95.8,
    'EnergyAware': 95,
    'Cobra': 95.6
}

# Color palette to ensure consistent colors across all plots
color_palette = {
    'Cobra': 'royalblue',
    'ECP': 'green',
    'EnergyAware': 'red',
    'Random': 'orange',
    'RoundRobin': 'darkviolet'
}

# Hatch patterns for the third plot
hatch_patterns = {
    'Cobra': 'xx',        # Cross pattern for Cobra
    'ECP': '||',          # Vertical lines for ECP
    'EnergyAware': '--',  # Horizontal lines for EnergyAware
    'Random': '//',       # Diag for Random
    'RoundRobin': '++'    # Crosshatch for RoundRobin
}

# Convert the data to a DataFrame
df_time = pd.DataFrame(task_time_data)
df_latency = pd.DataFrame(latency_data)  # Latency data for the fourth graph

# Function to plot UAV Battery Avg, UAV Available, Time Comparison, and Latency
def plot_uav_battery_and_availability_vs_tasks(csv_files):
    # Define markers for different models
    model_markers = {
        'Cobra': '^',  # Triangle for Cobra
        'ECP': 'd',  # Diamond for ECP
        'EnergyAware': 's',  # Square for EnergyAware
        'Random': 'x',  # Cross for Random
        'RoundRobin': 'o'  # Circle for RoundRobin
    }

    plt.figure(figsize=(10, 18))  # Adjusted figure height to fit all subplots

    # First subplot: UAV Battery Avg over Total Tasks
    plt.subplot(4, 1, 1)  # Four rows, one column, first plot
    for file in csv_files:
        model_name = file.split('_')[-1].split('.')[0]
        df = pd.read_csv(file)
        df['UAV Battery Avg'] = df['UAV Battery Avg'].apply(lambda x: 0 if x <= 5 else x)
        marker_style = model_markers.get(model_name.replace('TaskOffload', ''), 'o')
        plt.plot(df['Total Tasks'], df['UAV Battery Avg'], label=model_name.replace('TaskOffload', ''),
                 marker=marker_style, markevery=range(0, len(df['Total Tasks']), 10),
                 color=color_palette.get(model_name.replace('TaskOffload', ''), 'black'))  # Consistent colors
    plt.xlabel('Total Tasks Offload')
    plt.ylabel('Average UAV Battery (%)')
    plt.xlim(0, 1000)
    plt.xticks(range(0, 1001, 100))
    plt.ylim(0, 100)
    plt.grid(True)
    plt.legend()

    # Second subplot: UAV Available over Total Tasks
    plt.subplot(4, 1, 2)  # Four rows, one column, second plot
    for file in csv_files:
        model_name = file.split('_')[-1].split('.')[0]
        df = pd.read_csv(file)
        df['UAV Available'] = df['UAV Available'].apply(lambda x: max(x, 0))
        marker_style = model_markers.get(model_name.replace('TaskOffload', ''), 'o')
        plt.plot(df['Total Tasks'], df['UAV Available'], label=model_name.replace('TaskOffload', ''),
                 marker=marker_style, markevery=range(0, len(df['Total Tasks']), 5),
                 color=color_palette.get(model_name.replace('TaskOffload', ''), 'black'))  # Consistent colors
    plt.xlabel('Total Tasks Offload')
    plt.ylabel('Number of Available UAV')
    plt.xlim(500, 1000)
    plt.xticks(range(500, 1001, 100))
    plt.ylim(0)
    plt.grid(True)
    plt.legend()

    # Third subplot: Bar chart for time comparison with confidence intervals and hatches
    plt.subplot(4, 1, 3)  # Four rows, one column, third plot
    bar_width = 0.15
    tasks = df_time['Tasks']
    models = df_time.columns[1:]  # All models

    # Calculate the confidence intervals based on percentages provided
    error_bars = {model: (np.array(df_time[model]) * (100 - confidence_intervals[model]) / 100) for model in models}

    # Bar positions need to be shifted for better centering
    positions = np.arange(len(tasks))

    for i, model in enumerate(models):
        plt.bar(positions + i * bar_width, df_time[model], width=bar_width, label=model,
                yerr=error_bars[model], capsize=5, color=color_palette.get(model, 'black'),
                hatch=hatch_patterns.get(model, ''))  # Add hatches

    # Adjust xticks to be centered between the bars
    plt.xticks(positions + bar_width * (len(models) - 1) / 2, tasks)

    plt.xlabel('Tasks Offload Executed')
    plt.ylabel('Time Delay (S)')
    plt.ylim(0, 250)
    plt.grid(True)
    plt.legend()

    # Fourth subplot: Latency based on number of peers
    plt.subplot(4, 1, 4)  # Four rows, one column, fourth plot
    plt.plot(df_latency['NumPeer'], df_latency['TransactionTime'], label='Transaction Time', marker='o', color='blue')
    plt.plot(df_latency['NumPeer'], df_latency['ConsensusTime'], label='Consensus Time', marker='s', color='red')

    plt.xlabel('Number of Peers')
    plt.ylabel('Latency (S)')
    plt.xticks(range(0, 50, 5))  # Ensure proper spacing on the x-axis
    plt.yticks(range(0, 35, 5))  # Adjust y-axis ticks
    plt.grid(True)
    plt.legend()

    # Display the plots
    plt.tight_layout()
    plt.show()

    # Function to plot TotalTaskUAV (%) as a line plot for each model
def plot_totaltaskuav_lineplot(csv_files):
    plt.figure(figsize=(10, 6))  # Set up figure size for line plot

    # Read each CSV and extract relevant data
    for file in csv_files:
        model_name = file.split('_')[-1].split('.')[0]

        # Load the data
        df = pd.read_csv(file)

        # Group data every 50 tasks (as before) and calculate mean for grouping
        df_grouped = df.groupby(df['Total Tasks'] // 50).mean()  # Group every 50 tasks

        # Plot the line for each model
        plt.plot(df_grouped['Total Tasks'], df_grouped.iloc[:, -2], label=model_name.replace('TaskOffload', ''))

    # Customize the plot
    plt.xlabel('Total Tasks Offloaded')
    plt.ylabel('Proportion of task offload on the UAV (%)')
    plt.xlim(100, 1000)
    plt.xticks(range(100, 1001, 100))
    plt.ylim(0, 100)
    plt.grid(True)
    plt.legend()

    # Show the plot
    plt.tight_layout()
    plt.show()

# List of CSV files (assuming they are all in the current directory)
csv_files = glob.glob('graphe_result_*.csv')



# Call the function to plot all graphs
plot_uav_battery_and_availability_vs_tasks(csv_files)
# Call the function to plot the line graph
plot_totaltaskuav_lineplot(csv_files)
