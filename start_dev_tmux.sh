#!/bin/bash
SESSION="QuickGin"

# Start session and name the first window "backend"
tmux new-session -d -s $SESSION -n "backend"
tmux send-keys -t $SESSION:0 "cd ~/QuickGin && nvim" C-m

# Window 2: Web Frontend
tmux new-window -t $SESSION -n "frontend"
tmux send-keys -t $SESSION:1 "cd ~/QuickGin/web && nvim" C-m

# Window 3: Redis
tmux new-window -t $SESSION -n "redis"
tmux send-keys -t $SESSION:2 "redis-server" C-m

# Window 4: Postgres
tmux new-window -t $SESSION -n "pg"
tmux send-keys -t $SESSION:3 "pg_ctl -D \$PREFIX/var/lib/postgresql restart" C-m

# Window 5: Air (Go Hot Reload)
tmux new-window -t $SESSION -n "air"
tmux send-keys -t $SESSION:4 "cd ~/QuickGin && air" C-m

# Window 6: Web Dev (NPM)
tmux new-window -t $SESSION -n "web-run"
tmux send-keys -t $SESSION:5 "cd ~/QuickGin/web && npm run dev" C-m

# --- Custom Colors for Tabs ---
tmux set-window-option -t $SESSION:0 window-status-current-style bg=blue,fg=white   # Backend
tmux set-window-option -t $SESSION:1 window-status-current-style bg=cyan,fg=black  # Frontend
tmux set-window-option -t $SESSION:2 window-status-current-style bg=red,fg=white   # Redis
tmux set-window-option -t $SESSION:4 window-status-current-style bg=green,fg=black # Air

# Start on the Backend window
tmux select-window -t $SESSION:0
tmux attach-session -t $SESSION

