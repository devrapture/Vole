package scanner

type ImageAsset struct {
	AbsPath  string
	RelPath  string
	Basename string
	Used     bool
}

type ScanResult struct {
	ProjectPath  string        `json:"project_path"`
	AssetsDir    string        `json:"assets_dir"`
	TotalAssets  int           `json:"total_assets"`
	UsedAssets   int           `json:"used_assets"`
	UnusedAssets int           `json:"unused_assets"`
	Assets       []*ImageAsset `json:"assets"`
}

func (r *ScanResult) UnusedList() []*ImageAsset {
	var out []*ImageAsset
	for _, a := range r.Assets {
		if !a.Used {
			out = append(out, a)
		}
	}
	return out
}
