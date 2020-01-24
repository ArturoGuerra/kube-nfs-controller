package main

import (
    "os"
    "path"
    "errors"
    "strconv"
    "syscall"
    "github.com/golang/glog"
    "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
    "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
    Path = "path"
    Uid  = "uid"
    Gid  = "gid"
)


func (p *NfsProvisioner) Provision(options controller.ProvisionOptions) (*v1.PersistentVolume, error) {
    glog.Info("Provision called for volume: %s", options.PVName)

    fullPath, err := p.CreateOrGetShare(options)
    if err != nil {
        glog.Errorf("Failed to get or create NFS Share: %s Error: %s", options, err.Error())
        return nil, err
    }

    pv := &v1.PersistentVolume{
        ObjectMeta: metav1.ObjectMeta{
            Name: options.PVName,
        },
        Spec: v1.PersistentVolumeSpec{
            AccessModes: options.PVC.Spec.AccessModes,
            Capacity: v1.ResourceList{
                v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
            },
            PersistentVolumeReclaimPolicy: *options.StorageClass.ReclaimPolicy,
            PersistentVolumeSource: v1.PersistentVolumeSource{
                NFS: &v1.NFSVolumeSource{
                    Server: server,
                    Path: fullPath,
                    ReadOnly: false,
                },
            },
        },
    }

    return pv, nil
}

func (p *NfsProvisioner) CreateOrGetShare(options controller.ProvisionOptions) (string, error) {
    suffixPath := options.PVName

    struid, exists := options.StorageClass.Parameters[Uid]
    if !exists {
        struid = "65534"
    }

    uid, err := strconv.Atoi(struid)
    if err != nil {
        return "", err
    }

    strgid, exists := options.StorageClass.Parameters[Gid]
    if !exists {
        strgid = "65534"
    }

    gid, err := strconv.Atoi(strgid)
    if err != nil {
        return "", err
    }

    glog.Info("Fetching or Getting NFS Share")

    fullPath := path.Join(fakePath, suffixPath)

    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        syscall.Mkdir(fullPath, 0)
    }

    if err = syscall.Chown(fullPath, uid, gid); err != nil {
        return "", err
    }

    if err = syscall.Chmod(fullPath, 0775); err != nil {
        return "", err
    }

    return path.Join(mountPath, suffixPath), nil
}

func (p *NfsProvisioner) Delete(volume *v1.PersistentVolume) error {
    fullPath := path.Join(mountPath, volume.ObjectMeta.Name)

    if fullPath == mountPath {
        return errors.New("Sorry but you can't delete the root NFS Path")
    }

    if err := syscall.Rmdir(path.Join(fakePath, volume.ObjectMeta.Name)); err != nil {
        return err
    }

    return nil
}
