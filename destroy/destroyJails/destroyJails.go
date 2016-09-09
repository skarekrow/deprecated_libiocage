package destroyJails

import (
	"fmt"
	"github.com/iocage/libiocage/askUser"
	"github.com/iocage/libiocage/get/uclProp"
	"os"
	"os/exec"
	"strings"
)

func Args(Pool *string, force bool, jail, uuid, path, tag string) {
	var q bool

	if !force {
		qstr := fmt.Sprintf("\n"+strings.Repeat("-", 9)+"\nWARNING: \n"+
			strings.Repeat("-", 9)+"\n\n"+
			"This will destroy the jail: %s (%s)\n"+
			"Dataset: %s\n\n"+"Are you sure? y[N]: ", uuid, tag, path)
		q, _ = askUser.Args(qstr, true)
		fmt.Printf("\n")
	}

	if !q && !force {
		fmt.Printf("Command not confirmed.  No action taken.\n")
	} else {
		base, _ := uclProp.Args(path, "base", false)
		parent := *Pool + "/ioc/bases/" + base + "/root"

		// TODO: State check here
		_, err := exec.Command("/sbin/zfs", "destroy", "-Rv", parent+"@"+
			"jail_"+uuid).Output()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = exec.Command("/sbin/zfs", "destroy", *Pool+path).Output()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("%s successfully destroyed. (%s)\n", uuid, tag)
	}
}

/*
__destroy_jail () {
    local _name _answer _uuid_list _force _dataset _origin _fulluuid \
          _jail_path _state

    if [ -z "$1" ] ; then
        __die "missing UUID!"
    fi

    _name="$1"

    if [ "${_name}" = "ALL" ] ; then
        __die "please use iocage clean -j instead."
    fi

    _dataset="$(__find_jail ${_name})" || exit $?
    _fulluuid="$(__check_name ${_name} ${_dataset})"

    if [ -z ${_dataset} ] ; then
        __die "${_name} not found!"
    fi

    _origin="$(zfs get -H -o value origin ${_dataset})"
    _jail_path="$(__get_jail_prop mountpoint ${_fulluuid} ${_dataset})"
    _state="$(__is_running ${_fulluuid})"
    _jail_type="$(__get_jail_prop type ${_fulluuid} ${_dataset})"

    __check_children ${_fulluuid} ${_dataset}

    if [ "$?" -eq 1 ] ; then
        exit 1
    fi

    if [ "${_force}" -ne "1" ] ; then
            echo " "
            echo "  WARNING: this will destroy jail ${_fulluuid}"
            echo "  Dataset: ${_dataset}"
            echo " "
            echo -n "  Are you sure ? y[N]: "
            read _answer

        if [ "${_answer}" = "Y" -o "${_answer}" = "y" ] ; then
            if [ ! -z ${_state} ] ; then
                __die "cannot destroy ${_name} - jail is running!"
            fi

            __destroy_func ${_fulluuid} ${_dataset} ${_origin} ${_jail_type} ${_jail_path}
        else
            echo "  Command not confirmed.  No action taken."
        fi
    else
        if [ ! -z ${_state} ] ; then
            __stop_jail ${_fulluuid} ${_dataset}
        fi

        __destroy_func ${_fulluuid} ${_dataset} ${_origin} ${_jail_type} ${_jail_path}
    fi
}

__destroy_func () {
    local _fulluuid _dataset _origin _jail_type _base_inuse _jail_path

    _fulluuid="$1"
    _dataset="$2"
    _origin="$3"
    _jail_type="$4"
    _jail_path="$5"
    _base_inuse="$(zfs get -r -H -o value origin ${pool} | \
                 grep ${pool}/iocage/base > \
                 /dev/null 2>&1 ; echo $?)"

    echo "  Destroying: ${_fulluuid}"

    __unlink_tag ${_dataset}

    zfs destroy -fr ${_dataset}

    if [ "${_origin}" != "-" ] ; then
        echo "  Destroying clone origin: ${_origin}"
        zfs destroy -r ${_origin}

        if [ -d "${_jail_path}" ] ; then
            rm -rf ${_jail_path}
        fi
    fi

    if [ ${_jail_type} = "basejail" ] ; then
        if [ -d "${_jail_path}" ] ; then
            rm -rf ${_jail_path}
        fi
    fi

    if [ "${_base_inuse}" = "1" ] ; then
        zfs destroy -fr ${pool}/iocage/base > /dev/null 2>&1
    fi
}


__check_children () {
    local _dataset _jail_datasets _fs _origin _tag _fulluuid \
          _uuid _grepstring _printf

    _fulluuid="$1"
    _dataset="$2"
    _grepstring="jails$|base|releases|templates$|download|${pool}/iocage$"
    _grepstring="${_grepstring}|*./data"
    _jail_datasets=$(zfs list -d3 -rH -o name "${pool}/iocage" \
        | egrep -v "${_grepstring}")

    for _fs in ${_jail_datasets} ; do
        _origin=$(zfs get -H -o value origin "${_fs}" | cut -f1 -d@)

        if ! echo "${_fs}" | grep -q "/root" ; then
            # _fs is actually a dataset, so fulluuid needs to be faked.
            _uuid=$(__check_name "name" "${_fs}" 2> /dev/null)
            _tag=$(__get_jail_prop tag "${_uuid}" "${_fs}")
        fi

        if [ "${_origin}" = "${_dataset}/root" ] ; then
            _printf=" ERROR: jail has dependent clone, uuid: ${_uuid} (${_tag})"
            error="${error} $(printf "\n%s" "${_printf}")"
        fi
    done

    if [ ! -z "${error}" ] ; then
        export error
        return 1
    fi
}*/
