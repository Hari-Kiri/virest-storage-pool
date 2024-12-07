# ViRest Storage Pool
## Libvirt REST API Storage Pool Management
A storage pool is a quantity of storage set aside by an administrator, often a dedicated storage administrator, for use by virtual machines. Storage pools are divided into storage volumes either by the storage administrator or the system administrator, and the volumes are assigned to VMs as block devices.

For example, the storage administrator responsible for an NFS server creates a share to store virtual machines' data. The system administrator defines a pool on the virtualization host with the details of the share (e.g. nfs.example.com:/path/to/share should be mounted on /vm_data). When the pool is started, libvirt mounts the share on the specified directory, just as if the system administrator logged in and executed 'mount nfs.example.com:/path/to/share /vmdata'. If the pool is configured to autostart, libvirt ensures that the NFS share is mounted on the directory specified when libvirt is started.

Libvirt provides storage management on the physical host through storage pools and volumes and this software provides the interface using REST API.

## Needed Package to Running Executable
- qemu-kvm
- libvirt-daemon-system
- bridge-utils

## Needed Package for Development
- libvirt-dev
- gcc

## Add User to Libvirt Group & KVM Group
```shell
sudo adduser '<username>' libvirt
```
```shell
sudo adduser '<username>' kvm
```

## Known error:
- Libvirt Go Binding methods undefined, please enable "cgo" with command:
    ```shell
    export CGO_ENABLED=1
    ```
- [Can't access storage, file permission denied](https://ostechnix.com/solved-cannot-access-storage-file-permission-denied-error-in-kvm-libvirt/)

#### References
- https://libvirt.org/storage.html