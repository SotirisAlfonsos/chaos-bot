[![Build Status](https://travis-ci.org/SotirisAlfonsos/chaos-bot.svg)](https://travis-ci.org/SotirisAlfonsos/chaos-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/SotirisAlfonsos/chaos-bot)](https://goreportcard.com/report/github.com/SotirisAlfonsos/chaos-bot)
[![codebeat badge](https://codebeat.co/badges/9b7765b0-e40f-4534-8203-dde06d78dc07)](https://codebeat.co/projects/github-com-sotirisalfonsos-chaos-bot-master)
[![codecov](https://codecov.io/gh/SotirisAlfonsos/chaos-bot/branch/master/graph/badge.svg?token=ICGOMLDXRY)](https://codecov.io/gh/SotirisAlfonsos/chaos-bot)

# chaos-bot

Bots are responsible for the fault injections and should be run as services on the target servers with privileged system permissions (root).  
The [chaos master](https://github.com/SotirisAlfonsos/chaos-master) is responsible for controlling the bots. The master provides an api through which all fault injections can be orchestrated.

#### Fault injection types   
- &check; Docker outage: The selected Container is killed 
- &check; Services outage: The selected Service is killed
- &check; Cpu spikes: Create CPU spike based on percentage.
- &check; Server outage: The current server will be shut down
- &check; Network & package failures: Add network delays and package loss to the specific dev
- <i>(Coming soon)</i>
  - Memory spike injection
  - File descriptors spike injection
  - Kubernetes failure
