from syncloudlib import logger

from syncloud_platform.disks.hardware import has_unallocated_space_at_the_end
from syncloud_platform.disks.lsblk import Disk, Partition

logger.init(console=True)


def test_has_unallocated_space_at_the_end_low_percent():

    parted = '''BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;
1:100%:100%:0.00%:free;'''

    assert has_unallocated_space_at_the_end(parted) == False


def test_has_unallocated_space_at_the_end_high_percent():

    parted = '''BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;
1:100%:100%:12.34%:free;'''

    assert has_unallocated_space_at_the_end(parted) == True


def test_has_unallocated_space_at_the_end_no_free():

    parted = '''BYT;
/dev/sda:100%:scsi:512:512:msdos:ATA KINGSTON SV300S3:;
1:0.00%:0.00%:0.00%:free;
1:0.00%:1.67%:1.67%:::;
2:1.67%:100%:98.3%:ext4::;'''

    assert has_unallocated_space_at_the_end(parted) == False


