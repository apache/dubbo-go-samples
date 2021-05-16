package pkg

import (
    "context"
    "fmt"
    "math/rand"
    "strconv"

    "github.com/apache/dubbo-go/common/logger"

    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

type BasketballService struct{}

type Info struct {
    Name string
    Uid string
    Score int
}

var userMap = make(map[string]*Info, 0)

func (p *BasketballService) Online(ctx context.Context, uid string) (*pojo.Result, error) {
    logger.Infof("online: %s", uid)
    return &pojo.Result{Code: 0, Msg: "hello this is game provider"}, nil
}

func (p *BasketballService) Offline(ctx context.Context, uid string) (*pojo.Result, error) {
    logger.Infof("offline: %#s", uid)
    return &pojo.Result{Code: 0, Msg: "hello this is game provider"}, nil
}

func (p *BasketballService) Message(ctx context.Context, uid, data string) (*pojo.Result, error) {
    logger.Infof("message: %#v, %#v", uid, data)

    // auto reply the same message
    rsp, err := GateBasketball.Send(context.TODO(), uid, data)
    if err != nil {
        logger.Errorf("send fail: %#s", err.Error())
        return &pojo.Result{Code: 1, Msg: err.Error()}, err
    }

    fmt.Println("receive data from gate:", rsp)

    return &pojo.Result{Code: 0, Data: map[string]interface{}{"to": uid, "message": data + " from game provider"}}, nil
}

func (p *BasketballService) Login(ctx context.Context, data string) (*pojo.Result, error) {
    logger.Infof("message: %#v", data)
    var (
        info *Info
        ok bool
    )
    if info, ok = userMap[data]; !ok {
        fmt.Println("没有找到")
        info = &Info{}
        info.Uid = strconv.Itoa(rand.Int())
        info.Name = data
        userMap[info.Uid] = info
    }
    fmt.Println("找到了")
    return &pojo.Result{Code: 0, Data: map[string]interface{}{"to": info.Uid, "message": info.Name + ", your score is " + strconv.Itoa(info.Score)}}, nil
}

func (p *BasketballService) Score(ctx context.Context, uid string, score int) (*pojo.Result, error) {
    logger.Infof("message: %#v, %#v", uid, score)
    var (
        info *Info
        ok bool
    )
    if info, ok = userMap[uid]; !ok {
        logger.Error("user data not found")
    }
    info.Score += score
    return &pojo.Result{Code: 0, Data: map[string]interface{}{"to": info.Uid, "message": strconv.Itoa(info.Score)}}, nil
}

func (p *BasketballService) Reference() string {
    return "gameProvider.basketballService"
}