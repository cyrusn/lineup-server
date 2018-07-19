# A simple line up server for parents day

# TODO
## System crash when data grown up
- [x] should try to use database instead.
- [x] For the above issue, have to redesign the data structure.
- [x] Add a simple login feature for teacher and student helper
- [x] do not update the whole mapSchedule, should split it by form or even class
- [ ] update only for change related.
- [ ] use ab - Apache HTTP server benchmarking tool to test the benchmark in depoly state

## In front-end application
- [ ] add deboucing feature for font end application


## Forget websocket, use interval polling instead
- [ ] use websocket to ask client to update api, but remember to change the UI before update.
- [ ] websocket only trigger browser to update, browser should update through http Request
