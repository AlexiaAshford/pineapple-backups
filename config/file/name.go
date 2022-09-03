package config_file

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config/tool"
)

func NameSetting[T any](VolumeID, ChapOrder, ChapID T) string {
	return fmt.Sprintf("%v", VolumeID) + "-" +
		fmt.Sprintf("%05d", ChapOrder) + "-" + fmt.Sprintf("%d", ChapID) + ".txt"
}

func FileCacheName(index, chapter_index, ChapID any) string {
	index = tool.StrToInt(fmt.Sprintf("%d", index))
	return fmt.Sprintf("%05d", index) + "-" + fmt.Sprintf("%05d", chapter_index) + "-" +
		fmt.Sprintf("%v", ChapID) + ".txt"
}
