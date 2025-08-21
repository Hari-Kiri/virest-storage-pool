package storagePool

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolBuild"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func TestPoolBuildFromScratch(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	var incomingRequest poolBuild.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/build",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, libvirt.STORAGE_POOL_BUILD_NEW),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolBuild {
		test.Errorf("build pool from scratch test failed: %s", errorPoolBuild.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}

// Unsupported in libvirt version 10.0.0 and virsh (the libvirt command line interface) still not implement this method.
func TestPoolBuildRepairOrReinitialize(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := helper.helperGetStoragePoolFilesystemStruct(
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	var incomingRequest poolBuild.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/build",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, libvirt.STORAGE_POOL_BUILD_REPAIR),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolBuild {
		test.Errorf("repair or reinitilize pool test failed: %s", errorPoolBuild.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}

		if errorDeletePartition, isErrorDeletePartition := helper.helperDepleteDevicePartition(filesystemDiskDeviceValue); isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})
}

// Unsupported in libvirt version 10.0.0 and virsh (the libvirt command line interface) still not implement this method.
func TestPoolBuildExtendExistingPool(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := helper.helperGetStoragePoolFilesystemStruct(
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	var incomingRequest poolBuild.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/build",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, libvirt.STORAGE_POOL_BUILD_RESIZE),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolBuild {
		test.Errorf("extend existing pool test failed: %s", errorPoolBuild.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}

		if errorDeletePartition, isErrorDeletePartition := helper.helperDepleteDevicePartition(filesystemDiskDeviceValue); isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})
}

// Currently only filesystem pool accepts flags VIR_STORAGE_POOL_BUILD_NO_OVERWRITE (option 4).
func TestPoolBuildFilesystemNotOverwriteExistingPool(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := helper.helperGetStoragePoolFilesystemStruct(
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	var incomingRequest poolBuild.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/build",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, libvirt.STORAGE_POOL_BUILD_NO_OVERWRITE),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolBuild {
		test.Errorf("build not overwrite existing filesystem pool test failed: %s", errorPoolBuild.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}

		if errorDeletePartition, isErrorDeletePartition := helper.helperDepleteDevicePartition(filesystemDiskDeviceValue); isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})
}

// Currently only filesystem pool accepts flags VIR_STORAGE_POOL_BUILD_OVERWRITE (option 8).
func TestPoolBuildFilesystemOverwriteData(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := helper.helperGetStoragePoolFilesystemStruct(
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	var incomingRequest poolBuild.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/build",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, libvirt.STORAGE_POOL_BUILD_OVERWRITE),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolBuild {
		test.Errorf("build and overwrite existing filesystem pool test failed: %s", errorPoolBuild.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}

		if errorDeletePartition, isErrorDeletePartition := helper.helperDepleteDevicePartition(filesystemDiskDeviceValue); isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})
}

func (helper helperTest) helperTestPoolBuild(poolUuid string, option libvirt.StoragePoolBuildFlags) (virest.Error, bool) {
	helper.test.Helper()

	var incomingRequest poolBuild.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		helper.test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/build",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, option),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		helper.test.Fail()
		return errorHttpRequestPrecondition, isErrorHttpRequestPrecondition
	}
	defer poolConnection.Close()

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolBuild {
		helper.test.Fail()
		return errorPoolBuild, isErrorPoolBuild
	}

	return virest.Error{}, false
}
