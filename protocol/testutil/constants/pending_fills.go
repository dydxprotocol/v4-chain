package constants

import (
	"crypto/sha256"

	"github.com/dydxprotocol/v4/testutil/proto"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
)

var (
	// PendingFills.
	PendingFill_FullFill = clobtypes.PendingFill{
		Quantums:       10,
		Subticks:       15,
		TakerSide:      clobtypes.Order_SIDE_BUY,
		MakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20.Marshal())),
		TakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20.Marshal())),
		Type:           clobtypes.Trade,
	}
	PendingFill_PartialMakerFill1 = clobtypes.PendingFill{
		Quantums:       5,
		Subticks:       10,
		TakerSide:      clobtypes.Order_SIDE_SELL,
		MakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20.Marshal())),
		TakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.Marshal())),
		Type:           clobtypes.Trade,
	}
	PendingFill_PartialMakerFill2 = clobtypes.PendingFill{
		Quantums:       5,
		Subticks:       10,
		TakerSide:      clobtypes.Order_SIDE_SELL,
		MakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20.Marshal())),
		TakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15.Marshal())),
		Type:           clobtypes.Trade,
	}
	PendingFill_PartialTakerFill1 = clobtypes.PendingFill{
		Quantums:       5,
		Subticks:       10,
		TakerSide:      clobtypes.Order_SIDE_BUY,
		MakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.Marshal())),
		TakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20.Marshal())),
		Type:           clobtypes.Trade,
	}
	PendingFill_PartialTakerFill2 = clobtypes.PendingFill{
		Quantums:       5,
		Subticks:       10,
		TakerSide:      clobtypes.Order_SIDE_BUY,
		MakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15.Marshal())),
		TakerOrderHash: sha256.Sum256(proto.MustFirst(Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20.Marshal())),
		Type:           clobtypes.Trade,
	}
)
