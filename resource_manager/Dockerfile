FROM singaravelan21/mykubespray:v1.0.0
WORKDIR /
COPY . /app/resource-manager
ENV TMPDIR=/tmp
ENV GOMAXPROCS=8
EXPOSE 9090
ENTRYPOINT ["/app/resource-manager/main"]


#TASK [ : prep_kubeadm_images | Create kubeadm config] ***************************************
#fatal: [node2]: FAILED! => {"changed": false, "checksum": "2c6bbdc767a8200145ff13a21d5954a82a98d2b1", "msg": "Destination directory /etc/kubernetes does not exist"}
#fatal: [node1]: FAILED! => {"changed": false, "checksum": "2c6bbdc767a8200145ff13a21d5954a82a98d2b1", "msg": "Destination directory /etc/kubernetes does not exist"}

