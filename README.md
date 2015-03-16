Phi Accrual Failure Detection
=============================

Godoc: http://godoc.org/github.com/vektra/go-failure

Based on the algorithm described in http://ddg.jaist.ac.jp/pub/HDY+04.pdf.

Given a window of samples calculated from the time heartbeats are recieved,
this is able to calculate a suspicion level that the node is still operating
correctly.

This package uses the modified math version described in
https://issues.apache.org/jira/browse/CASSANDRA-2597 which keeps everything
failure simple.
