package report

type Reports []Report

func (r Reports) AddReport(report Report) (Reports) {
    var tmp []Report
    var newlen = len(r)
    tmp, r = r, make(Reports, newlen)
    copy(r, tmp)
    r = append(r, report)
    return r
}

func (r Reports) Names() []string {
    n := make([]string, len(r))
    for _, r := range r {
        n = append(n, r.FileName)
    }
    return n
}

func (r Reports) Len() int {
    return len(r)
}

func (r Reports) Less(i, j int) bool {
    return r[i].FileName > r[j].FileName;
}

func (r Reports) Swap(i, j int) {
    r[i], r[j] = r[j], r[i]
}
