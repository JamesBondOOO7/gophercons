<h1> Commands </h1>

<ul>
<li>
Run a specific test and dump the trace

```
go test -run TestCopyIt -trace=copy_trace.out
```
</li>
<li>
Analyse the trace

```
go tool trace copy_trace.out
```
</li>
</ul>