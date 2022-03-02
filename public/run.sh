echo "Running script..."
./kill_active_server
nohup ./serve >some.log 2>&1 </dev/null &
