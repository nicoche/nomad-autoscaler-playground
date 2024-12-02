# There can only be a single job definition per file.
# Create a job with ID and Name 'example'
job "redis" {
	# Run the job in the global region, which is the default.
	# region = "global"

	# Specify the datacenters within the region this job can run in.
	datacenters = ["dc1"]

	# Configure the job to do rolling updates
	update {
		# Stagger updates every 10 seconds
		stagger = "10s"

		# Update a single task at a time
		max_parallel = 1
	}

	# Create a 'cache' group. Each task in the group will be
	# scheduled onto the same machine.
	group "cache" {
		# Control the number of instances of this group.
		# Defaults to 1
		# count = 1

		# Configure the restart policy for the task group. If not provided, a
		# default is used based on the job type.
		restart {
			# The number of attempts to run the job within the specified interval.
			attempts = 2
			interval = "1m"

			# A delay between a task failing and a restart occurring.
			delay = "10s"

			# Mode controls what happens when a task has restarted "attempts"
			# times within the interval. "delay" mode delays the next restart
			# till the next interval. "fail" mode does not restart the task if
			# "attempts" has been hit within the interval.
			mode = "fail"
		}

    scaling {
      enabled = true
      min     = 0
      max     = 10

        policy {
          check = [
            {
              "cpu_usage": [
                {
                  "source": "prometheus",
                  "query": "cpu",
                  "group": "cpu",
                  "strategy": [
                    {
                      "target-value": [
                        {
                          "target": 50,
                          "max_scale_down": 1,
                          "threshold": 0.1
                        }
                      ]
                    }
                  ]
                }
              ],
              "scaletozero": [
                {
                  "query": "scaletozero",
                  "source": "prometheus",
                  "group": "cpu",
                  "strategy": [
                    {
                      "threshold": [
                        {
                          "lower_bound" = 0,
                          "upper_bound" = 0.9,
                          "delta" = -10000,
                        }
                      ]
                    }
                  ]
                }
              ]
            },
            {
              "reqpersecond": [
                {
                  "source": "prometheus",
                  "query": "reqpersecond",
                  "group": "reqpersecond",
                  "strategy": [
                    {
                      "pass-through": [
                        {
                          "max_scale_down": 1,
                        }
                      ]
                    }
                  ]
                }
              ],
              "scaletozero": [
                {
                  "query": "scaletozero",
                  "source": "prometheus",
                  "group": "reqpersecond",
                  "strategy": [
                    {
                      "threshold": [
                        {
                          "lower_bound" = 0,
                          "upper_bound" = 0.9,
                          "delta" = -500,
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
          cooldown = "15s"
          evaluation_interval = "15s"
        }
    }

		# Define a task to run
		task "redis" {
			# Use Docker to run the task.
			driver = "docker"

			# Configure Docker driver with the image
			config {
				image = "redis:latest"
				port_map {
					db = 6379
				}
			}

			resources {
				network {
					port "db" {
					}
				}
			}
		}
	}
}
