#Analyze Sunsetting plugin
![sg_analyze_light](https://user-images.githubusercontent.com/2936828/48772107-0b305300-eccc-11e8-8c72-4bcbd737226b.png)

[![Coverage Status](https://coveralls.io/repos/github/supergiant/analyze-plugin-sunsetting/badge.svg?branch=master)](https://coveralls.io/github/supergiant/analyze-plugin-sunsetting?branch=master)(https://coveralls.io/github/supergiant/analyze?branch=master)
[![Build Status](https://travis-ci.com/supergiant/analyze-plugin-sunsetting.svg?branch=master)](https://travis-ci.com/supergiant/analyze-plugin-sunsetting)
[![License Apache 2](https://img.shields.io/badge/License-Apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/supergiant/analyze-plugin-sunsetting)](https://goreportcard.com/report/github.com/supergiant/analyze-plugin-sunsetting)

Analyze sunsetting plugin is a micro-service which is made as Analyze plugin. It's workflow consists of next steps:
1. Collect CPU/RAM requests from pods and allocatable resources from nodes in a cluster.
2. Calculate the most effective way of pods placement on nodes, trying to calculate how to pack pods in way where some nodes could be removed as totally non-utilized (with no pods scheduled on them).
3. Final result is pushed to Analyze and is rendered as table, for instance: 

|Region/Zone	|Instance ID		|RAM requested (GIB)|RAM not requested (GIB)|Total RAM (GIB)|Price per day (USD)|Recommended to sunset|
|---------------|:-----------------:|:-----------------:|:---------------------:|:-------------:|:-----------------:|:-------------------:|
|ap-southeast-2b|i-0af2668d717b9cd14|0.252				|8.013					|8.265			|3					|Yes					|
|ap-southeast-2a|i-0e17d3572f2c09bbe|0.184				|8.082					|8.265			|3					|Yes					|
|ap-southeast-2a|i-0116cd0c83868bb33|0.91				|7.355					|8.265			|3					|No						|
|ap-southeast-2a|i-0b68adff84f64a3dd|0.597				|7.668					|8.265			|3					|No						|


####TL;DR 
**[Get started here](https://supergiant.readme.io/docs/node-sunsetting-plugin)**

##Implementations notes
1. Currently only AWS is supported
2. Backend of a plugin is written in Golang and frontend is written in Angular and is packed as set of web components which are built-in in binary.
3. On bootsrtaping time (right after installation) it fetches prices from AWS pricing API. This is single time operation.
4. Works only when ProviderID is set.


##Roadmap:
1. Current implementation do not assume that cluster shall survive single node failure.  
2. Current implementation do not counts that some pods can't be rescheduled based on:  
		* Taints and Tolerations  
		* Node affinity  
		* Inter-pod affinity/anti-affinity  
		* PV Claims  
		* etc...  
