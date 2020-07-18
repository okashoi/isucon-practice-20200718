.PHONY: up down logs bench restart kataribe

up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

bench: workload=1
bench: init=./init.sh
bench:
	sudo /opt/isucon3-mod/bench/bench benchmark --workload ${workload} --init ${init}

SLOW_LOG=/var/lib/mysql/mysql-slow.log
restart:
	sudo supervisorctl stop isucon-app && $(MAKE) -C app build && sudo supervisorctl start isucon-app
	sudo bash -c "echo '' > /var/log/nginx/access.log && echo '' > /var/log/nginx/error.log" && sudo systemctl restart nginx.service
	sudo bash -c "echo '' > $(SLOW_LOG)" && sudo systemctl restart mysqld.service

kataribe:
	cat /var/log/nginx/access.log | kataribe

torch:
	TIMESTAMP=`date "+%Y%m%d%H%M%S"` && go-torch --url http://localhost:6060/debug/pprof/profile -p > /opt/isucon3-mod/app/src/public/torch-$TIMESTAMP.svg && ln -sf /opt/isucon3-mod/app/src/public/torch-$TIMESTAMP.svg /opt/isucon3-mod/app/src/public/torch.svg
